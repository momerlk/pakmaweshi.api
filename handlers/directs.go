package handlers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"pakmaweshi.api/internal"
)

func (a *App) Directs(w http.ResponseWriter , r *http.Request){
	claims , ok := Verify(w , r)
	if !ok { return }

	userId := claims["user_id"]

	directs , err := internal.Get[internal.Direct](r.Context() , &a.Database , "directs" , bson.M{
		"user_id" : userId,
	})
	if err != nil {
		a.ServerError(w , "Directs" , err)
		return
	}

	json.NewEncoder(w).Encode(directs)
}