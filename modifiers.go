package handler

import (
	"fmt"
)

func ModifyRequestBody(body string, contentType string, conf FluentLoggerConfig) string {
	return modifyBody(body, contentType, conf.Response.allowedContentTypes, conf.Response.bodyLimit)
}

func ModifyResponseBody(body string, contentType string, conf FluentLoggerConfig) string {
	return modifyBody(body, contentType, conf.Request.allowedContentTypes, conf.Request.bodyLimit)
}

func modifyBody(body string, contentType string, contentTypes map[string]struct{}, limit int) string {
	runes := []rune(body)

	if _, ok := contentTypes[contentType]; ok {
		if len(runes) > limit {
			return string(runes[:limit])
		}
		return body
	}
	return fmt.Sprintf("Content-Type's \"%s\" body not allowed to log", contentType)
}
