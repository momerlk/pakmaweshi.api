package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"pakmaweshi.api/handlers"
	"pakmaweshi.api/internal"

	"github.com/rs/cors"
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

	ws := &internal.WebSocket{};
	ws.Init(func(conn *internal.WSConnection , data []byte) {
		log.Println("onRead called")
	});

	mux.Handle("/test" , ws)
	mux.HandleFunc("/upload" , app.UploadFile);
	mux.HandleFunc("/file", app.DownloadFile)
	mux.HandleFunc("/signUp" , app.SignUp)
	mux.HandleFunc("/signIn" , app.SignIn)
	mux.HandleFunc("/feed" , app.Feed)
	mux.HandleFunc("/post" , app.CreatePost);

	
	handler := cors.New(cors.Options{
		AllowedOrigins : []string{
			"http://localhost:19000",
			"http://192.168.18.6:19000",
		},
		AllowCredentials : true,


		Debug : true,
	}).Handler(mux)

	PORT := os.Getenv("PORT")
	http.ListenAndServe("localhost:" + PORT , handler)
}
