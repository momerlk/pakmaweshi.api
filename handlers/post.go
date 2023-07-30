package handlers

import (
	"context"
	"encoding/json"
	"log"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"pakmaweshi.api/internal"

	"net/http"
)


const productsColl = "products"


func (a *App) CreatePost(w http.ResponseWriter , r *http.Request){
	var claims jwt.MapClaims
	var ok bool
	if claims , ok = a.Verify(w , r); !ok {
		return
	}




	userId := claims["user_id"].(string)


	var data internal.Product

	data.Id = internal.GenerateId();
	data.UserId = userId;
	var user internal.User
	ok , err := a.Database.Get(r.Context() , "users" , bson.M{"id" : userId} , &user);
	if !ok || err != nil {
		log.Printf("no user found, err = %v" , err)
		return
	}

	data.Username = user.Username
	data.Avatar = user.Avatar


	if r.Method != http.MethodPost {
		http.Error(w , http.StatusText(http.StatusMethodNotAllowed) , http.StatusMethodNotAllowed)
		return
	}

	data.Id = uuid.NewString()

	err = json.NewDecoder(r.Body).Decode(&data)	
	if err != nil {
		log.Println(err)
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return 
	}


	err = a.Database.Store(context.TODO() , productsColl , data)
	if err != nil {
		http.Error(w , err.Error() , http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Successfully posted data"))
	return
}