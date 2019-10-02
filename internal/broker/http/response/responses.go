package response

import (
	"net/http"

	"github.com/dipress/crmifc/internal/validation"
	"github.com/pkg/errors"
)

// easyjson -all responses.go

var (
	badRequestBody = messageResponse{
		Message: "bad request",
	}

	internalServerErrorBody = messageResponse{
		Message: "internal server error",
	}

	notFoundBody = messageResponse{
		Message: "not found",
	}

	unauthorizedBody = messageResponse{
		Message: "unauthorized",
	}
)

type messageResponse struct {
	Message string `json:"message"`
}

// BadRequestResponse returns status bad request.
func BadRequestResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusBadRequest)

	data, err := badRequestBody.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// InternalServerErrorResponse returns internal server error.
func InternalServerErrorResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusInternalServerError)

	data, err := internalServerErrorBody.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}

// NotFoundResponse returns not found response.
func NotFoundResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)

	data, err := notFoundBody.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}
	return nil
}

// UnauthorizedResponse returns unauthorized response.
func UnauthorizedResponse(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusUnauthorized)

	data, err := unauthorizedBody.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}
	return nil
}

type validationResponse struct {
	Message string            `json:"message"`
	Errors  validation.Errors `json:"errors"`
}

// UnprocessabeEntityResponse returns unprocessabe entity response.
func UnprocessabeEntityResponse(w http.ResponseWriter, ers validation.Errors) error {
	w.WriteHeader(http.StatusUnprocessableEntity)

	ver := validationResponse{
		Message: ers.Error(),
		Errors:  ers,
	}

	data, err := ver.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "write response")
	}

	return nil
}
