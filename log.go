package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const emptyQty = 0
const minTokenParts = 2

type LogData struct {
	start              time.Time
	path               string
	fulle_path         string
	clientIP           string
	host               string
	requestMethod      string
	requestHeaders     http.Header
	requestBody        string
	responseStatusCode int
	responseHeaders    http.Header
	rawResponseBody    *bytes.Buffer
	responseBody       string
}

type LogWriter struct {
	gin.ResponseWriter
	writer  io.Writer
	logData LogData
}

func (lw LogWriter) Write(b []byte) (int, error) {
	return lw.writer.Write(b)
}

func (lw LogWriter) GetHeaderValue(key string) string {
	return lw.logData.requestHeaders.Get(key)
}

func bytesToString(bytes []byte) string {
	return string(bytes)
}

func (lw *LogWriter) SetRequestBody(c *gin.Context, conf FluentLoggerConfig) {
	lw.logData.requestBody = ModifyRequestBody(c, conf)
}

func (lw *LogWriter) SetResponseBody(c *gin.Context, conf FluentLoggerConfig) {
	lw.logData.responseHeaders = c.Writer.Header()
	lw.logData.responseStatusCode = c.Writer.Status()
	lw.logData.responseBody = ModifyResponseBody(c, lw.logData.rawResponseBody, conf)
}

func (lw *LogWriter) MakeLogData(conf FluentLoggerConfig) map[string]interface{} {
	data := lw.logData
	finish := time.Now()
	maskedRequestHeaders := MaskRequestHeaders(makeHeaders(data.requestHeaders), conf.Mask.Request)
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
		"request.body":         data.requestBody,
		"response.status_code": fmt.Sprintf("%v", data.responseStatusCode),
		"response.headers":     createKeyValuePairs(maskedResponseHeader),
		"response.body":        data.responseBody,
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

func AddHeader(c *gin.Context, header_key string, header_value string) {
	c.Request.Header.Set(header_key, header_value)
}

func NewLogWriter(c *gin.Context) (*LogWriter, error) {
	var log bytes.Buffer

	full_path := ""
	raw := c.Request.URL.RawQuery

	if raw != "" {
		full_path = c.Request.URL.Path + "?" + raw
	}

	newLogWriter := &LogWriter{
		ResponseWriter: c.Writer,
		writer:         io.MultiWriter(c.Writer, &log),
		logData: LogData{
			start:           time.Now(),
			path:            c.Request.URL.Path,
			fulle_path:      full_path,
			clientIP:        c.ClientIP(),
			host:            c.Request.Host,
			requestHeaders:  c.Request.Header,
			requestMethod:   c.Request.Method,
			rawResponseBody: &log,
		},
	}

	c.Writer = newLogWriter

	return newLogWriter, nil
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
