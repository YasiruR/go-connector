package middleware

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/tryfix/log"
	"io"
	"net/http"
)

var logger pkg.Log

func init() {
	logger = log.Constructor.Log(
		log.WithColors(true),
		log.WithLevel(log.TRACE),
		log.WithFilePath(true),
		log.WithSkipFrameCount(4),
	)
}

type Parser struct {
	typ api.HandlerType
}

// to compose error messages based on handler
func NewParser(typ api.HandlerType) *Parser {
	return &Parser{typ: typ}
}

func ParseRequest(r *http.Request, val any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body.Close()
		return errors.InvalidRequestBody(``, err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &val); err != nil {
		return errors.UnmarshalError(``, err)
	}

	return nil
}

func WriteAck(w http.ResponseWriter, data any, statusCode int) {
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			WriteError(w, errors.WriteAckFailed(``, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		if _, err = w.Write(body); err != nil {
			WriteError(w, errors.WriteAckFailed(``, err), http.StatusInternalServerError)
			return
		}

		return
	}

	w.WriteHeader(statusCode)
	return
}

func WriteError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	logger.Error(err)
}
