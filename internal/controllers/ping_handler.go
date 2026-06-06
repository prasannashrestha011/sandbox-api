package controllers

import (
	"main/internal/response"
	"net/http"
)

type PingerController struct {
}

func NewPingerController() *PingerController {
	return &PingerController{}
}
func (c *PingerController) Ping(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, r, http.StatusOK, "pong", map[string]string{"response": "pong"}, nil)
}
