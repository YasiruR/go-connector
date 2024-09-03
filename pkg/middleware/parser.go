package middleware

import (
	"encoding/json"
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/dsp"
	"github.com/YasiruR/connector/domain/errors/external"
	"github.com/YasiruR/connector/domain/pkg"
	pkgLog "github.com/tryfix/log"
	"io"
	"net/http"
)

var log pkg.Log

func init() {
	log = pkgLog.Constructor.Log(
		pkgLog.WithColors(true),
		pkgLog.WithLevel(pkgLog.TRACE),
		pkgLog.WithFilePath(true),
		pkgLog.WithSkipFrameCount(3),
	)
}

func ParseRequest(r *http.Request, val any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		r.Body.Close()
		return core.ReadBodyFailed(err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &val); err != nil {
		return core.UnmarshalError(err)
	}

	return nil
}

func WriteAck(w http.ResponseWriter, data any, statusCode int) error {
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return core.WriteAckFailed(err)
		}

		w.WriteHeader(statusCode)
		if _, err = w.Write(body); err != nil {
			return core.WriteAckFailed(err)
		}

		return nil
	}

	w.WriteHeader(statusCode)
	return nil
}

func WriteError(w http.ResponseWriter, err error, statusCode int) {
	var ge external.GatewayError
	var ce dsp.CatalogError
	var ne dsp.NegotiationError
	var te dsp.TransferError

	var data []byte
	var tmpErr error

	switch {
	case defaultErr.As(err, &ge):
		data, tmpErr = json.Marshal(ge)
		fmt.Println("GEEEE")
	case defaultErr.As(err, &ne):
		fmt.Println("NEEEE")
		data, tmpErr = json.Marshal(ne.Body)
	case defaultErr.As(err, &te):
		fmt.Println("TEEEE")
		data, tmpErr = json.Marshal(te.Body)
	case defaultErr.As(err, &ce):
		fmt.Println("CEEEE")
		data, tmpErr = json.Marshal(ce.Body)
	}

	fmt.Println("DAATAA: ", string(data))

	w.WriteHeader(statusCode)
	if tmpErr == nil {
		if _, tmpErr = w.Write(data); tmpErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = fmt.Errorf("%w AND %s", err, tmpErr)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err = fmt.Errorf("%w AND %s", err, tmpErr)
	}

	log.Error(err)
}
