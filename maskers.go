package handler

import (
	"encoding/json"
	"fmt"
	"strings"
)

func MaskRequestHeaders(data map[string]string, conf map[string][]string) map[string]string {
	headers, ok := conf["request.headers"]
	if !ok {
		return data
	}

	return MaskHeaders(data, headers)
}

func MaskResponseHeaders(data map[string]string, conf map[string][]string) map[string]string {
	headers, ok := conf["response.headers"]
	if !ok {
		return data
	}

	return MaskHeaders(data, headers)
}

func MaskHeaders(data map[string]string, conf []string) map[string]string {
	for _, header := range conf {
		value, ok := data[header]
		if !ok {
			continue
		}

		splitted := strings.Split(value, " ")
		if len(splitted) == 2 {
			if splitted[0] == "Bearer" || splitted[0] == "Token" {
				forJoin := []string{splitted[0], maskFormat(splitted[1])}
				data[header] = strings.Join(forJoin, " ")
			}
		} else {
			joined := strings.Join(splitted, " ")
			data[header] = maskFormat(joined)
		}

	}

	return data
}

func MaskRequestBody(body string, conf map[string][]string) string {
	keys, ok := conf["request.body"]
	if !ok {
		return body
	}

	return MaskBody(body, keys)
}

func MaskResponseBody(body string, conf map[string][]string) string {
	keys, ok := conf["response.body"]
	if !ok {
		return body
	}

	return MaskBody(body, keys)
}

func MaskBody(body string, conf []string) string {
	jsonBody, err := toJSON(body)
	if err != nil {
		return body
	}

	for _, key := range conf {
		value, ok := jsonBody[key]
		if !ok {
			continue
		}
		jsonBody[key] = maskFormat(fmt.Sprint(value))
	}

	jsonString, err := toString(jsonBody)

	if err != nil {
		return body
	}

	return jsonString
}

func maskFormat(s string) string {
	var formatted string

	runes := []rune(s)

	if len(s) >= 11 {
		formatted = fmt.Sprint(string(runes[:4]), "...", string(runes[len(runes)-4:]))
	} else {
		formatted = strings.Repeat("*", len(s))
	}

	return formatted
}

func toJSON(data string) (map[string]interface{}, error) {
	bufferSingleMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &bufferSingleMap)

	return bufferSingleMap, err
}

func toString(s map[string]interface{}) (string, error) {
	jsonString, err := json.Marshal(s)

	return string(jsonString), err
}
