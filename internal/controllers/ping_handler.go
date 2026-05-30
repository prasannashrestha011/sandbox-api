package controllers

import "net/http"

type PingerController struct {
}

func NewPingerController() *PingerController {
	return &PingerController{}
}
func (c *PingerController) Ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"response": "pong"})
}
