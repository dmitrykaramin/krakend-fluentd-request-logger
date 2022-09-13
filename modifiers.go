package handler

import (
	"fmt"
)

func ModifyResponseBody(body string, contentType string, conf FluentLoggerConfig) string {
	runes := []rune(body)

	if _, ok := conf.Response.allowedContentTypes[contentType]; ok {
		if len(runes) > conf.Response.bodyLimit {
			return string(runes[:conf.Response.bodyLimit])
		}
		return body
	}
	return fmt.Sprintf("Content-Type's \"%s\" body not allowed to log", contentType)
}
