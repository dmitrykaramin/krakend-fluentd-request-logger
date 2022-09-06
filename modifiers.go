package handler

import (
	"fmt"
)

func ModifyResponseBody(data LogData, conf FluentLoggerConfig) string {
	body := data.responseBody.String()
	contentType := data.requestHeaders.Get("Content-Type")

	if _, ok := conf.Response.allowedContentTypes[contentType]; ok {
		if len(body) > conf.Response.bodyLimit {
			return body[:conf.Response.bodyLimit]
		}
		return body
	}
	return fmt.Sprintf("Content-Type's \"%s\" body not allowed to log", contentType)
}
