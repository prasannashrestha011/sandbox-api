package response

import (
	"encoding/json"
	"main/internal/domain"
	"main/internal/dto"
	"net/http"
)

func WriteJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	msg string,
	data any,
	meta any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := dto.Response{
		Success: status >= 200 && status < 300,
		Message: msg,
		Data:    data,
		Error:   nil,
	}

	_ = json.NewEncoder(w).Encode(res)
}

func WriteJSONError(
	w http.ResponseWriter,
	r *http.Request,
	err *domain.AppError,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPCode)

	res := dto.Response{
		Success: false,
		Message: err.Message,
		Error: &dto.ErrorBody{
			Code:    string(err.Code),
			Details: err.ErrorDetails,
		},
	}

	_ = json.NewEncoder(w).Encode(res)
}
