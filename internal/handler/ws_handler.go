package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"sespima_api/internal/ws"
	pkgjwt "sespima_api/pkg/jwt"
	"sespima_api/pkg/response"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:    func(r *http.Request) bool { return true },
}

type WSHandler struct {
	hub       *ws.Hub
	jwtSecret string
}

func NewWSHandler(hub *ws.Hub, jwtSecret string) *WSHandler {
	return &WSHandler{hub: hub, jwtSecret: jwtSecret}
}

var wsViewerRoles = map[string]bool{
	"korsis": true, "patun": true, "pimpinan": true, "operator": true, "kabag_bindik": true,
}

// StreamLocations upgrades an HTTP GET to WebSocket for live location viewing.
// Auth is passed as ?token=<jwt> since WS clients can't set headers after the upgrade.
func (h *WSHandler) StreamLocations(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		response.Unauthorized(c, "token diperlukan")
		return
	}
	claims, err := pkgjwt.ParseAccessToken(h.jwtSecret, tokenStr)
	if err != nil {
		response.Unauthorized(c, "token tidak valid")
		return
	}
	if !wsViewerRoles[claims.Role] {
		response.Forbidden(c, "akses ditolak")
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(h.hub, conn)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
