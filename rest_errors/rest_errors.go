package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	message string        `json:"message"`
	status  int           `json:"status"`
	error   string        `json:"error"`
	causes  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf(
		"message: %s - status: %d - error: %s - causes: [ %v ]",
		e.message, e.status, e.error, e.causes,
	)
}

func (e restErr) Message() string {
	return e.message
}

func (e restErr) Status() int {
	return e.status
}

func (e restErr) Causes() []interface{} {
	return e.causes
}

// func NewError(msg string) error {
// 	return errors.New(msg)
// }

func NewRestError(
	message string, status int, err string, causes []interface{},
) RestErr {
	return restErr{
		message: message,
		status:  status,
		error:   err,
		causes:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(msg string) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusBadRequest,
		error:   "bad_request",
	}
}

func NewNotFoundError(msg string) RestErr {
	return restErr{
		message: msg,
		status:  http.StatusNotFound,
		error:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		message: message,
		status:  http.StatusUnauthorized,
		error:   "unauthorized",
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	ret := restErr{
		message: msg,
		status:  http.StatusInternalServerError,
		error:   "internal_server_error",
	}
	if err != nil {
		ret.causes = append(ret.causes, err.Error())
	}

	return ret
}
