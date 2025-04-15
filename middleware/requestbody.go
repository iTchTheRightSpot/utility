package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/iTchTheRightSpot/utility/utils"
	"io"
	"net/http"
)

type RequestBodyMiddleware[T any] struct {
	Logger    utils.ILogger
	Validator *validator.Validate
}

func (dep *RequestBodyMiddleware[T]) RequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			dep.Logger.Error(r.Context(), "request body is nil")
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		var payload T

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			dep.Logger.Error(r.Context(), err.Error())
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		if err := dep.Validator.Struct(payload); err != nil {
			dep.Logger.Error(r.Context(), err.Error())
			utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid request body"})
			return
		}

		by, err := json.Marshal(payload)
		if err != nil {
			dep.Logger.Error(r.Context(), err.Error())
			utils.ErrorResponse(w, &utils.ServerError{})
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(by))

		next.ServeHTTP(w, r)
	})
}