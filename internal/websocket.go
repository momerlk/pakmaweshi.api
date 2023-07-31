package internal

import (
	"log"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v5"


	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

)

func Verify(w http.ResponseWriter , r *http.Request) (jwt.MapClaims , bool) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return nil , false
	}

	var tokenClaims jwt.MapClaims

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key := Getenv("JWT_KEY")
		return []byte(key) , nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil , false
	}

	// Check if the token is valid and not expired
	if claims , ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		http.Error(w, "Invalid token (expired)", http.StatusUnauthorized)
		return nil , false
	} else {
		tokenClaims = claims
	}


	return tokenClaims , true
}


// represents a websocket endpoint
type WebSocket struct {
	epoller 		*Epoll
}

func (socket *WebSocket) Init(onRead func(conn *WSConnection , data []byte)){
	// Start epoll
	var err error
	epoller, err := MkEpoll()
	if err != nil {
		panic(err)
	}

	socket.epoller = epoller;

	go socket.Start(onRead)
}

func (socket *WebSocket) Send(userId string , data []byte) (bool , error) {
	conn , ok := socket.epoller.Writers[userId]
	if !ok {
		return false , nil
	}

	conn.Lock.Lock()
	defer conn.Lock.Unlock()

	err := wsutil.WriteClientBinary(conn.Conn , data)
	if err != nil {
		return false , err
	}

	return true , nil
}

func (socket *WebSocket) Start(onRead func(conn *WSConnection , data []byte)){
	for {
		connections, err := socket.epoller.Wait()
		
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if msg, _, err := wsutil.ReadClientData(conn.NetConn); err != nil {
				if err := socket.epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.NetConn.Close()
			} else {
				log.Printf("msg : %s" , string(msg))
				onRead(conn , msg)
				// This is commented out since in demo usage, stdout is showing messages sent from > 1M connections at very high rate
				//log.Printf("msg: %s", string(msg))
			}
		}
	}
}

func (socket *WebSocket) ServeHTTP(w http.ResponseWriter , r *http.Request){

	claims , ok := Verify(w , r)
	if !ok {
		return
	}

	userId := claims["user_id"].(string)


	// Upgrade connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := socket.epoller.Add(&WSConnection{NetConn : conn, UserId: userId}); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()
	}
}