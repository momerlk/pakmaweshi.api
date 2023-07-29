package main

import (
	"pakmaweshi.api/handlers"
	"pakmaweshi.api/internal"
	"net/http"
)

type HttpHandler func (w http.ResponseWriter , r *http.Request)

func POST(w http.ResponseWriter , r *http.Request , handler HttpHandler) HttpHandler{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w , http.StatusText(http.StatusMethodNotAllowed) , http.StatusMethodNotAllowed)
			return
		}
		handler(w , r)
	}	
}

func GET(w http.ResponseWriter , r *http.Request , handler HttpHandler) HttpHandler{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w , http.StatusText(http.StatusMethodNotAllowed) , http.StatusMethodNotAllowed)
			return
		}
		handler(w , r)
	}	
}

func main(){
	mux := http.NewServeMux();

	db := internal.Database{}
	db.Init()

	app := handlers.App{
		Database: db,
	}

	mux.HandleFunc("/post" , app.CreatePost);
	mux.HandleFunc("/upload" , app.UploadFile);
	mux.HandleFunc("/file", app.DownloadFile)
	mux.HandleFunc("/signUp" , app.SignUp)
	mux.HandleFunc("/signIn" , app.SignIn)
	mux.HandleFunc("/private" , app.HandlePrivate)


	http.ListenAndServe(":8080" , mux)
}
