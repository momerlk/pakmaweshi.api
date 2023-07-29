package handlers

import (
	"log"
	"time"

	"encoding/json"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"pakmaweshi.api/internal"

	"net/http"
)

const usersColl = "users"


func (a *App) SignUp(w http.ResponseWriter , r *http.Request){
	var body internal.User
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.ServerError(w , "Sign Up" , err)
		return
	}

	body.Id = uuid.NewString()

	hashed , err := internal.HashAndSalt([]byte(body.Password))
	if err != nil {
		a.ServerError(w , "Sign Up" , err)
		return
	}
	body.Password = hashed
	

	a.Database.Store(r.Context() , usersColl , body)

	w.Write([]byte("successfully registered user"))
}


type SignInBody struct {
	UsernameEmail	 			string 			`json:"username_email" bson:"username_email"` 			// username or email
	Password 					string 			`json:"password" bson:"password"`
}
func (a *App) SignIn(w http.ResponseWriter , r *http.Request){
	var body SignInBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.ServerError(w , "Sign In" , err)
		return
	}



	var user internal.User
	ok , err := a.Database.Get(r.Context() , usersColl , bson.D{{"username" , body.UsernameEmail}} , &user)
	if err != nil {
		a.ServerError(w , "Sign In a.Database.Get()" , err)
		return
	}
	if !ok {
		ok , err := a.Database.Get(r.Context() , usersColl , bson.D{{"email" , body.UsernameEmail}} , &user)
		if !ok {
			a.ClientError(w , http.StatusUnauthorized)
			return
		}
		if err != nil {
			a.ServerError(w , "Sign In" , err)
			return
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password) , []byte(body.Password)) == nil {
		log.Println("user is authenticated")
		// User is authenticated
		secret := internal.Getenv("JWT_KEY")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256 , jwt.MapClaims{
			"user_id" : user.Id,
			"session_id" : internal.GenerateId(),
			"exp" : time.Now().Add(4 * time.Hour).Unix(),
		})

		tokenString , err := token.SignedString([]byte(secret))
		if err != nil {
			a.ServerError(w , "Sign In" , err)
			return
		}

		err = json.NewEncoder(w).Encode(bson.M{"token" : tokenString})
		if err != nil {
			a.ServerError(w , "Sign In" , err)
			return
		}

	} else {
		a.ClientError(w , http.StatusUnauthorized)
		return
	}

}

func (a *App) Verify(w http.ResponseWriter , r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token missing", http.StatusUnauthorized)
		return false
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key := internal.Getenv("JWT_KEY")
		return []byte(key) , nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return false
	}

	// Check if the token is valid and not expired
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		http.Error(w, "Invalid token (expired)", http.StatusUnauthorized)
		return false
	}

	return true
}

