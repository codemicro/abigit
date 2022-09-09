package util

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type RichError struct {
	Status int
	Reason string
	Detail any
}

func NewRichError(status int, reason string, detail any) error {
	return &RichError{
		Status: status,
		Reason: reason,
		Detail: detail,
	}
}

func NewRichErrorFromFiberError(err *fiber.Error, detail any) error {
	return NewRichError(err.Code, err.Message, detail)
}

func (r *RichError) Error() string {
	return fmt.Sprintf("handler error, %d: %s", r.Status, r.Reason)
}

func (r *RichError) AsJSON() ([]byte, error) {
	info := map[string]any{
		"status":  "error",
		"message": r.Reason,
	}

	if r.Detail != nil {
		info["detail"] = r.Detail
	}

	return json.Marshal(info)
}

func JSONErrorHandler(ctx *fiber.Ctx, err error) error {
	causeErr := errors.Cause(err) // if an error has been wrapped with a stacktrace, it'll mess up this matching below
	var re *RichError
	if e, ok := causeErr.(*fiber.Error); ok {
		re = NewRichErrorFromFiberError(e, nil).(*RichError)
	} else if e, ok := causeErr.(*RichError); ok {
		re = e
	} else {
		log.Error().Stack().Err(err).Str("location", "fiber error handler").Str("route", ctx.OriginalURL()).Send()
		re = NewRichErrorFromFiberError(fiber.ErrInternalServerError, nil).(*RichError)
	}
	jsonBytes, err := re.AsJSON()
	if err != nil {
		jsonBytes = []byte(`{"status":"error","message":"Internal Server Error","detail":"unable to produce detailed description"}`)
		log.Error().Err(err).Str("location", "fiber error handler").Msg("unable to produce error response")
	}
	ctx.Status(re.Status)
	ctx.Type("json")
	return ctx.Send(jsonBytes)
}
