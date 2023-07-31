package handlers

import (
	"context"
	"encoding/json"

	"pakmaweshi.api/internal"
)

const directsColl =  "directs"

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
		return err
	}

	err = a.Database.Store(context.TODO() , directsColl , direct)
	if err != nil {
		return err
	}

	return err
}


