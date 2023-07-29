package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"pakmaweshi.api/internal"

	"net/http"
)


const productsColl = "products"


func (a *App) CreatePost(w http.ResponseWriter , r *http.Request){
	var data internal.Product

	if r.Method != "POST" {
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