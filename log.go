package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

const emptyQty = 0
const minTokenParts = 2

type LogData struct {
	start              time.Time
	path               string
	clientIP           string
	host               string
	requestMethod      string
	requestHeaders     http.Header
	requestBody        []byte
	responseStatusCode int
	responseHeaders    http.Header
	responseBody       *bytes.Buffer
}

type LogWriter struct {
	gin.ResponseWriter
	writer  io.Writer
	logData LogData
}

func (lw LogWriter) Write(b []byte) (int, error) {
	return lw.writer.Write(b)
}

func (lw *LogWriter) MakeLogData(conf FluentLoggerConfig) map[string]interface{} {
	data := lw.logData
	finish := time.Now()
	contentType := data.responseHeaders.Get("Content-Type")
	maskedRequestBody := MaskRequestBody(string(data.requestBody), conf.Mask.Request)
	maskedRequestHeaders := MaskRequestHeaders(makeHeaders(data.requestHeaders), conf.Mask.Request)
	maskedResponseBody := MaskResponseBody(data.responseBody.String(), conf.Mask.Response)
	maskedResponseHeader := MaskResponseHeaders(makeHeaders(data.responseHeaders), conf.Mask.Response)

	return map[string]interface{}{
		"start":                fmt.Sprintf("%v", data.start),
		"finish":               fmt.Sprintf("%v", finish),
		"path":                 data.path,
		"latency":              fmt.Sprintf("%v", finish.Sub(data.start)),
		"client_ip":            data.clientIP,
		"host":                 data.host,
		"request.method":       data.requestMethod,
		"request.headers":      createKeyValuePairs(maskedRequestHeaders),
		"request.body":         maskedRequestBody,
		"response.status_code": fmt.Sprintf("%v", data.responseStatusCode),
		"response.headers":     createKeyValuePairs(maskedResponseHeader),
		"response.body":        ModifyResponseBody(maskedResponseBody, contentType, conf),
	}
}

func AddJwtData(data map[string]interface{}, claimsToAdd map[string]struct{}, header string) error {
	if header == "" || len(claimsToAdd) <= emptyQty {
		return nil
	}

	splitHeader := strings.Split(header, " ")

	if len(splitHeader) < minTokenParts {
		return errors.New("wrong authorization header format")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(splitHeader[1], jwt.MapClaims{})
	if err != nil {
		return err
	}

	for k, v := range token.Claims.(jwt.MapClaims) {
		if _, ok := claimsToAdd[k]; ok {
			data[k] = v
		}
	}

	return nil
}

func NewLogWriter(c *gin.Context) (*LogWriter, error) {
	var log bytes.Buffer

	bodyToRead, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	if raw != "" {
		path = path + "?" + raw
	}

	newLogWriter := &LogWriter{
		ResponseWriter: c.Writer,
		writer:         io.MultiWriter(c.Writer, &log),
		logData: LogData{
			start:          time.Now(),
			path:           path,
			clientIP:       c.ClientIP(),
			host:           c.Request.Host,
			requestBody:    bodyToRead,
			requestHeaders: c.Request.Header,
			requestMethod:  c.Request.Method,
			responseBody:   &log,
		},
	}

	c.Writer = newLogWriter
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyToRead))

	return newLogWriter, nil
}

func (lw *LogWriter) CompleteLogData(c *gin.Context) {
	lw.logData.responseHeaders = c.Writer.Header()
	lw.logData.responseStatusCode = c.Writer.Status()
}

func makeHeaders(header http.Header) map[string]string {
	headers := make(map[string]string)
	for k, vs := range header {
		headers[k] = strings.Join(vs, ", ")
	}

	return headers
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
