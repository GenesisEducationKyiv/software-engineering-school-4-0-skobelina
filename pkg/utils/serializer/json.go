package serializer

import (
	"encoding/json"
	"io"
	"net/http"

	errors "github.com/skobelina/currency_converter/pkg/utils/errors"
)

// swagger:model message
type JsonMessage struct {
	Message string `json:"message"`
}

func SendJSON(w http.ResponseWriter, status int, object interface{}) error {
	// set headers
	SetCorsHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	// write response
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func SendNoContent(w http.ResponseWriter) error {
	// set headers
	SetCorsHeaders(w)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func SendError(w http.ResponseWriter, err errors.Error) error {
	return SendJSON(w, errors.ParseStatusCode(err),
		&map[string]string{"type": errors.ParseErrorType(err), "message": err.Error()},
	)
}

func SetCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func ParseJsonBody(body io.ReadCloser, value interface{}) error {
	if err := json.NewDecoder(body).Decode(value); err != nil && !errors.Is(err, io.EOF) {
		return errors.NewBadRequestError(err)
	}
	return nil
}
