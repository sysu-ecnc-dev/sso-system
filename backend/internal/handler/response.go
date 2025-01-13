package handler

import (
	"net/http"

	"github.com/sysu-ecnc-dev/sso-system/backend/internal/utils"
)

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func successResponse(w http.ResponseWriter, message string, data any) {
	utils.WriteJSON(w, http.StatusOK, response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}
