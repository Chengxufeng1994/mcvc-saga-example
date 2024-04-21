package model

import (
	"encoding/json"
	"strings"
)

const maxErrorLength = 1024

type AppError struct {
	Id            string `json:"id"`
	Message       string `json:"message"`              // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"`       // Internal error string to help the developer
	RequestId     string `json:"request_id,omitempty"` // The RequestId that's also set in the header
	Where         string `json:"-"`                    // The function where it happened in the form of Struct.Func
	params        map[string]any
	wrapped       error
}

func NewAppError(where string, id string, params map[string]any, details string) *AppError {
	ap := &AppError{
		Id:            id,
		params:        params,
		Message:       id,
		Where:         where,
		DetailedError: details,
	}
	return ap
}

func (er *AppError) Error() string {
	var sb strings.Builder

	// render the error information
	if er.Where != "" {
		sb.WriteString(er.Where)
		sb.WriteString(": ")
	}

	// only render the detailed error when it's present
	if er.DetailedError != "" {
		sb.WriteString(er.DetailedError)
	}

	// render the wrapped error
	err := er.wrapped
	if err != nil {
		sb.WriteString(", ")
		sb.WriteString(err.Error())
	}

	res := sb.String()
	if len(res) > maxErrorLength {
		res = res[:maxErrorLength] + "..."
	}
	return res
}

func (er *AppError) ToJSON() string {
	// turn the wrapped error into a detailed message
	detailed := er.DetailedError
	defer func() {
		er.DetailedError = detailed
	}()

	er.wrappedToDetailed()

	b, _ := json.Marshal(er)
	return string(b)
}

func (er *AppError) wrappedToDetailed() {
	if er.wrapped == nil {
		return
	}

	if er.DetailedError != "" {
		er.DetailedError += ", "
	}

	er.DetailedError += er.wrapped.Error()
}

func (er *AppError) Unwrap() error {
	return er.wrapped
}

func (er *AppError) Wrap(err error) *AppError {
	er.wrapped = err
	return er
}

func (er *AppError) WipeDetailed() {
	er.wrapped = nil
	er.DetailedError = ""
}
