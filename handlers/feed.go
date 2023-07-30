package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"pakmaweshi.api/internal"
)

func (a *App) Feed(w http.ResponseWriter , r *http.Request){
	claims , ok := a.Verify(w , r)	
	if !ok {
		a.ClientError(w , http.StatusUnauthorized)
		return
	}
	_ = claims["user_id"]

	var posts internal.Product
	ok , err := a.Database.Get(r.Context() , "products" , map[string]string{"" : ""} , &posts);
	if err != nil {
		a.ServerError(w , "Feed" , err)
		return
	} 
	if !ok {
		http.Error(w , "No Posts" , http.StatusExpectationFailed)
		return
	}

	log.Println("posts =" , posts)

	json.NewEncoder(w).Encode(posts)

}