package internal


type Product struct {
	Id				string 					`json:"id" bson:"id"`

	Name 			string					`json:"name" bson:"name"`		
	Description		string					`json:"description" bson:"description"`								
	Price 			string					`json:"price" bson:"price"`
	Location 		string					`json:"location" bson:"location"`
	Contact 		string					`json:"contact" bson:"contact"`
	
	Avatar 			string					`json:"avatar" bson:"avatar"`
	Username 		string					`json:"username" bson:"username"`  // username
	UserId			string					`json:"user_id" bson:"user_id"`	 // user id

	Images 			[]string				`json:"images" bson:"images"` // url of the images
}

type HiddenProduct struct {

	Name 			string					`json:"name" bson:"name"`		
	Description		string					`json:"description" bson:"description"`								
	Price 			string					`json:"price" bson:"price"`
	Location 		string					`json:"location" bson:"location"`
	Contact 		string					`json:"contact" bson:"contact"`
	
	Avatar 			string					`json:"avatar" bson:"avatar"`
	Username 		string					`json:"username" bson:"username"`  // username

	Images 			[]string				`json:"images" bson:"images"` // url of the images
}

type User struct {
	Id 				string 					`json:"id" bson:"id"` // user id 

	Avatar			string 					`json:"avatar" bson:"avatar"` // url of the avatar image file

	Name 			string					`json:"name" bson:"name"` // full name
	Number 			string 					`json:"number" bson:"number"` // phone number only +92
	Username		string 					`json:"username" bson:"username"` // username

	Email			string 					`json:"email" bson:"email"` // email
	Password		string 					`json:"password" bson:"password"`	 // password
}

type Direct struct {
	Id 				string 					`json:"id" bson:"id"` 					// message id

	Sender 			string 					`json:"sender" bson:"sender"` 			// sender's user id
	Receiver		string 					`json:"receiver" bson:"receiver"` 		// receiver's user id
	Received 		bool					`json:"received" bson:"received"` 		// whether the message has been received or not

	Content 		string 					`json:"content" bson:"content"` 	   // text content of the message
	Attachment 		string 					`json:"attachment" bson:"attachment"` // file id of the attachment
}

type RenderedDirect struct {
	Content 		string 					`json:"content" bson:"content"`
	TimeSent 		string 					`json:"time_sent" bson:"time_sent"`
	Sent 			bool					`json:"sent"`
}

type RenderedChat struct {
	Name 			string					`json:"name"`
	Username 		string					`json:"username"`
	Avatar 			string					`json:"avatar"`
	Messages 		[]RenderedDirect 		`json:"messages"`
}