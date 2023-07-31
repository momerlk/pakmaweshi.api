package handlers

import (

	"encoding/json"

	_ "github.com/gobwas/ws/wsutil"
	"pakmaweshi.api/internal"
)


func (a *App) Direct(conn *internal.WSConnection , data []byte) (err error) {
	var direct internal.Direct

	err = json.Unmarshal(data , &direct)

	
	


	return err
}