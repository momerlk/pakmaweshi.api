package handlers

import (
	"context"
	"encoding/json"
	"log"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"pakmaweshi.api/internal"

	"net/http"
)


const productsColl = "products"


func (a *App) CreatePost(w http.ResponseWriter , r *http.Request){
	var token *jwt.Token
	var ok bool
	if token , ok = a.Verify(w , r); !ok {
		return
	}

	userId := token.Header["user_id"].(string)

	var data internal.Product

	data.Id = internal.GenerateId();
	data.UserId = userId;
	var user internal.User
	a.Database.Get(r.Context() , "users" , map[string]string{"user_id" : userId} , &user);

	data.Username = user.Username


	if r.Method != http.MethodPost {
		http.Error(w , http.StatusText(http.StatusMethodNotAllowed) , http.StatusMethodNotAllowed)
		return
	}

	data.Id = uuid.NewString()

	err := json.NewDecoder(r.Body).Decode(&data)	
	if err != nil {
		log.Println(err)
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return 
	}

	log.Println(data)

	err = a.Database.Store(context.TODO() , productsColl , data)
	if err != nil {
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Successfully posted data"))
	return
}