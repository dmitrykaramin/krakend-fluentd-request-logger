package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const emptyQty = 0
const minTokenParts = 2

type LogData struct {
	start              time.Time
	path               string
	clientIp           string
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

func (lw *LogWriter) MakeLogData() map[string]string {
	data := lw.logData
	finish := time.Now()

	return map[string]string{
		"start":                fmt.Sprintf("%v", data.start),
		"finish":               fmt.Sprintf("%v", finish),
		"path":                 data.path,
		"latency":              fmt.Sprintf("%v", finish.Sub(data.start)),
		"client_ip":            data.clientIp,
		"host":                 data.host,
		"request.method":       data.requestMethod,
		"request.headers":      makeHeaders(data.requestHeaders),
		"request.body":         string(data.requestBody),
		"response.status_code": fmt.Sprintf("%v", data.responseStatusCode),
		"response.headers":     makeHeaders(data.requestHeaders),
		"response.body":        fmt.Sprintf("%v", data.responseBody.String()),
	}
}

func AddJwtData(data map[string]string, claimsToAdd map[string]struct{}, header string) error {
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
			data[k] = fmt.Sprintf("%v", v)
		}
	}

	return nil
}

func NewLogWriter(c *gin.Context) (*LogWriter, error) {
	var log bytes.Buffer

	bodyToRead, err := ioutil.ReadAll(c.Request.Body)
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
			clientIp:       c.ClientIP(),
			host:           c.Request.Host,
			requestBody:    bodyToRead,
			requestHeaders: c.Request.Header,
			requestMethod:  c.Request.Method,
			responseBody:   &log,
		},
	}

	c.Writer = newLogWriter
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyToRead))

	return newLogWriter, nil
}

func (lw *LogWriter) CompleteLogData(c *gin.Context) {
	lw.logData.responseHeaders = c.Writer.Header()
	lw.logData.responseStatusCode = c.Writer.Status()
}

func makeHeaders(header http.Header) string {
	headers := make(map[string]string)
	for k, vs := range header {
		headers[k] = strings.Join(vs, ", ")
	}

	return createKeyValuePairs(headers)
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
