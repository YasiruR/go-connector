package middleware

import (
	"encoding/json"
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
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
		return readBodyFailed(err)
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &val); err != nil {
		return unmarshalError(err)
	}

	return nil
}

func WriteAck(w http.ResponseWriter, data any, statusCode int) error {
	if data != nil {
		body, err := json.Marshal(data)
		if err != nil {
			return writeAckFailed(err)
		}

		w.WriteHeader(statusCode)
		if _, err = w.Write(body); err != nil {
			return writeAckFailed(err)
		}

		return nil
	}

	w.WriteHeader(statusCode)
	return nil
}

func WriteError(w http.ResponseWriter, err error, statusCode int) {
	var clientErr errors.ClientError
	var catErr errors.CatalogError
	var negErr errors.NegotiationError
	var trnErr errors.TransferError

	var res any
	var tmpErr error
	switch {
	case defaultErr.As(err, &clientErr):
		res = clientErr
	case defaultErr.As(err, &negErr):
		res = parseNegotiationErr(negErr)
	case defaultErr.As(err, &trnErr):
		res = parseTransferErr(trnErr)
	case defaultErr.As(err, &catErr):
		res = parseCatalogErr(catErr)
	}

	data, tmpErr := json.Marshal(res)
	if tmpErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(fmt.Errorf("%w AND %s", err, tmpErr))
		return
	}

	w.WriteHeader(statusCode)
	if _, tmpErr = w.Write(data); tmpErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = fmt.Errorf("%w AND %s", err, tmpErr)
	}

	log.Error(err)
}

func parseCatalogErr(ce errors.CatalogError) catalog.Error {
	return catalog.Error{
		Context:    control_plane.Context,
		Type:       catalog.MsgTypeError,
		DspaceCode: ce.Code,
		DspaceReason: []struct {
			Value    string `json:"@value"`
			Language string `json:"@language"`
		}{{Value: ce.Message, Language: `en`}},
	}
}

func parseNegotiationErr(ne errors.NegotiationError) negotiation.Error {
	return negotiation.Error{
		Ctx:     control_plane.Context,
		Type:    negotiation.MsgTypeError,
		ProvPId: ne.ProvPid,
		ConsPId: ne.ConsPid,
		Code:    ne.Code,
		Reason:  []interface{}{ne.Message},
		Desc:    nil,
	}
}

func parseTransferErr(te errors.TransferError) transfer.Error {
	return transfer.Error{
		Ctx:     control_plane.Context,
		Type:    transfer.MsgTypeError,
		ProvPId: te.ProvPid,
		ConsPId: te.ConsPid,
		Code:    te.Code,
		Reason:  []interface{}{te.Message},
	}
}
