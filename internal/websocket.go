package internal

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"net/http"
)

type WebSocket struct {
	epoller 		*Epoll
}

func (socket *WebSocket) Init(onRead func(conn net.Conn , data []byte)){
	// Start epoll
	var err error
	epoller, err := MkEpoll()
	if err != nil {
		panic(err)
	}

	socket.epoller = epoller;

	go socket.Start(onRead)
}


func (socket *WebSocket) Start(onRead func(conn net.Conn , data []byte)){
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
			if msg, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := socket.epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
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
	// Upgrade connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := socket.epoller.Add(conn); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()
	}
}