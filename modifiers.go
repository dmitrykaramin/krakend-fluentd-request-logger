package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

func ModifyRequestBody(c *gin.Context, conf FluentLoggerConfig) string {
	requestContentType := c.Request.Header.Get("Content-Type")
	if ok := checkContentType(requestContentType, conf); !ok {
		return fmt.Sprintf("Request Content-Type's \"%s\" body not allowed to log", requestContentType)
	}

	if ok := checkContentLength(c.Request.ContentLength, conf); !ok {
		return fmt.Sprintf("Content too long  \"%s\" ", c.Request.ContentLength)
	}

	bodyToRead, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("Error reading body: \"%s\"", err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyToRead))

	return string(bodyToRead)
}

func ModifyResponseBody(c *gin.Context, responseBody *bytes.Buffer, conf FluentLoggerConfig) string {
	requestContentType := c.Writer.Header().Get("Content-Type")
	if ok := checkContentType(requestContentType, conf); !ok {
		return fmt.Sprintf("Response Content-Type's \"%s\" body not allowed to log", requestContentType)
	}

	var readResponseBody []byte
	var err error

	contentLength := int64(responseBody.Len())
	if ok := checkContentLength(contentLength, conf); !ok {
		bodyReader := io.LimitReader(responseBody, conf.Response.bodyLimit)
		readResponseBody, err = io.ReadAll(bodyReader)
		if err != nil {
			return fmt.Sprintf("Error reading response body: \"%s\"", err.Error())
		}

	} else {
		readResponseBody, err = io.ReadAll(responseBody)
		if err != nil {
			return fmt.Sprintf("Error reading response body: \"%s\"", err.Error())
		}
	}

	return string(readResponseBody)
}

func checkContentType(contentType string, conf FluentLoggerConfig) bool {
	if _, ok := conf.Request.allowedContentTypes[contentType]; ok {
		return true
	} else {
		return false
	}
}

func checkContentLength(contentLength int64, conf FluentLoggerConfig) bool {
	if contentLength < conf.Request.bodyLimit {
		return true
	} else {
		return false
	}
}
