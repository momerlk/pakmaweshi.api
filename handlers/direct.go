package handlers

import (
	"context"
	"encoding/json"

	"pakmaweshi.api/internal"
)

const directsColl =  "directs"

type WSDirectMeta struct {
	Received 			bool 			`json:"received" bson:"received"`
}

// websocket handler for directs
func (a *App) WSDirect(ws *internal.WebSocket , conn *internal.WSConnection , data []byte) (err error) {
	var direct internal.Direct

	err = json.Unmarshal(data , &direct)
	if err != nil {
		return err
	}

	ok , err := ws.Send(direct.Receiver , data)
	if err != nil {
		return err
	}

	if ok {
		direct.Received = true
	} else {
		direct.Received = false
	}

	err = a.Database.Store(context.TODO() , directsColl , direct)
	if err != nil {
		return err
	}

	meta , err := json.Marshal(WSDirectMeta{
		Received: direct.Received,
	})
	if err != nil {
		return err
	}

	ws.Send(direct.Sender , meta)

	return err
}


