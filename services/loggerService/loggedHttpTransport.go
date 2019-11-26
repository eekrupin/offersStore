package loggerService

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	maxResponseBodyLogLength = 32768
)

type LoggedRoundTripper struct {
	Ctx                  context.Context
	IsResponseLogEnabled bool
}

func getBodyAsString(body *io.ReadCloser, maxLength int) string {
	if *body == nil {
		return ""
	}

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(*body)
	*body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if maxLength > 0 && len(bodyBytes) > maxLength {
		bodyBytes = bodyBytes[0:maxLength]
	}

	return string(bodyBytes)
}

func (tripper LoggedRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	startTime := time.Now()
	requestBody := getBodyAsString(&req.Body, 0)
	GetMainLogger().Info(tripper.Ctx,
		"Sending HTTP request",
		"method", req.Method,
		"host", req.Host,
		"path", req.URL.Path,
		"requestBody", requestBody,
	)

	res, err = http.DefaultTransport.RoundTrip(req)

	duration := time.Since(startTime)
	if err != nil {
		GetMainLogger().Warn(tripper.Ctx, "Received HTTP response (error)",
			"method", req.Method,
			"host", req.Host,
			"path", req.URL.Path,
			"requestBody", requestBody,
			"error", err.Error(),
			"duration", duration,
		)
	} else if tripper.IsResponseLogEnabled {
		GetMainLogger().Info(tripper.Ctx, "Received HTTP response",
			"method", req.Method,
			"host", req.Host,
			"path", req.URL.Path,
			"requestBody", requestBody,
			"responseBody", getBodyAsString(&res.Body, maxResponseBodyLogLength),
			"duration", duration,
		)
	}

	return res, err
}
