package handlers

import (
	"net"

	"encoding/json"

	"github.com/gobwas/ws/wsutil"
	"pakmaweshi.api/internal"
)


func (a *App) Direct(conn net.Conn , data []byte) (err error) {
	var direct internal.Direct

	err = json.Unmarshal(data , &direct)

	return err
}