package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent" //nolint:goimports
	"github.com/luraproject/lura/v2/logging"
)

type MaskConfig struct {
	Request  map[string][]string
	Response map[string][]string
}

type BodyLoggerConfig struct {
	bodyLimit           int64
	allowedContentTypes map[string]struct{}
}

type FluentLoggerConfig struct {
	FluentTag    string
	FluentConfig fluent.Config
	Skip         map[string]struct{}
	logger       logging.Logger
	JWTClaims    map[string]struct{}
	FlatHeaders  map[string]struct{}
	Response     BodyLoggerConfig
	Request      BodyLoggerConfig
	Mask         MaskConfig
}

func printOutConfigError(key string, err error) {
	message := "used default value for '%s' fluentd config error: %v \n"
	printOutError(key, err, message)
}

func printOutError(key string, err error, message string) {
	m := fmt.Sprintf("krakend-fluentd-request-logger: %v", message)

	fmt.Printf(m, key, err)
}

func ConvertToString(key string, cfg map[string]interface{}) string {
	err := errors.New("no value found")
	value, ok := cfg[key]
	if ok {
		return fmt.Sprintf("%v", value)
	}
	printOutConfigError(key, err)

	return ""
}

func ConvertToInt(key string, cfg map[string]interface{}) int {
	err := errors.New("no value found")
	value, ok := cfg[key]
	valueString := fmt.Sprintf("%v", value)
	if ok {
		p, err := strconv.Atoi(valueString)
		if err == nil {
			return p
		}
	}
	printOutConfigError(key, err)

	return 0
}

func ConvertToBool(key string, cfg map[string]interface{}) bool {
	err := errors.New("no value found")
	value, ok := cfg[key]
	if ok {
		p, ok := value.(bool)
		if ok {
			return p
		}
	}
	printOutConfigError(key, err)

	return false
}

func (f *FluentLoggerConfig) setFluentHost(cfg map[string]interface{}) {
	f.FluentConfig.FluentHost = ConvertToString("fluent_host", cfg)
}

func (f *FluentLoggerConfig) setFluentPort(cfg map[string]interface{}) {
	f.FluentConfig.FluentPort = ConvertToInt("fluent_port", cfg)
}

func (f *FluentLoggerConfig) setFluentNetwork(cfg map[string]interface{}) {
	f.FluentConfig.FluentNetwork = ConvertToString("fluent_network", cfg)
}

func (f *FluentLoggerConfig) setFluentSocketPath(cfg map[string]interface{}) {
	f.FluentConfig.FluentSocketPath = ConvertToString("fluent_socket_path", cfg)
}

func (f *FluentLoggerConfig) setFluentTimeout(cfg map[string]interface{}) {
	f.FluentConfig.Timeout = time.Duration(ConvertToInt("timeout", cfg))
}

func (f *FluentLoggerConfig) setFluentWriteTimeout(cfg map[string]interface{}) {
	f.FluentConfig.WriteTimeout = time.Duration(ConvertToInt("write_timeout", cfg))
}

func (f *FluentLoggerConfig) setBufferLimit(cfg map[string]interface{}) {
	f.FluentConfig.BufferLimit = ConvertToInt("buffer_limit", cfg)
}

func (f *FluentLoggerConfig) setRetryWait(cfg map[string]interface{}) {
	f.FluentConfig.RetryWait = ConvertToInt("retry_wait", cfg)
}

func (f *FluentLoggerConfig) setMaxRetry(cfg map[string]interface{}) {
	f.FluentConfig.MaxRetry = ConvertToInt("max_retry", cfg)
}

func (f *FluentLoggerConfig) setMaxRetryWait(cfg map[string]interface{}) {
	f.FluentConfig.MaxRetryWait = ConvertToInt("max_retry_wait", cfg)
}

func (f *FluentLoggerConfig) setTagPrefix(cfg map[string]interface{}) {
	f.FluentConfig.TagPrefix = ConvertToString("tag_prefix", cfg)
}

func (f *FluentLoggerConfig) setAsync(cfg map[string]interface{}) {
	f.FluentConfig.Async = ConvertToBool("async", cfg)
}

func (f *FluentLoggerConfig) setForceStopAsyncSend(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("force_stop_async_send", cfg)
}

func (f *FluentLoggerConfig) setSubSecondPrecision(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("sub_second_precision", cfg)
}

func (f *FluentLoggerConfig) setRequestAck(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("request_ack", cfg)
}

func (f *FluentLoggerConfig) setTag(cfg map[string]interface{}) {
	f.FluentTag = ConvertToString("fluent_tag", cfg)
}

func (f *FluentLoggerConfig) SetFluentConfig(cfg map[string]interface{}) error {
	fluentConfig, ok := cfg["fluent_config"]
	if !ok {
		return errors.New("no 'fluent_config' key found. using default fluent config")
	}

	fluentConfigMap, ok := fluentConfig.(map[string]interface{})
	if !ok {
		return errors.New("can't convert config to right type. using default fluent config")
	}

	f.setFluentHost(fluentConfigMap)
	f.setFluentPort(fluentConfigMap)
	f.setFluentNetwork(fluentConfigMap)
	f.setFluentSocketPath(fluentConfigMap)
	f.setFluentTimeout(fluentConfigMap)
	f.setFluentWriteTimeout(fluentConfigMap)
	f.setBufferLimit(fluentConfigMap)
	f.setRetryWait(fluentConfigMap)
	f.setMaxRetry(fluentConfigMap)
	f.setMaxRetryWait(fluentConfigMap)
	f.setTagPrefix(fluentConfigMap)
	f.setAsync(fluentConfigMap)
	f.setForceStopAsyncSend(fluentConfigMap)
	f.setSubSecondPrecision(fluentConfigMap)
	f.setRequestAck(fluentConfigMap)
	f.setTag(fluentConfigMap)

	return nil
}

func (f *FluentLoggerConfig) SetSkipConfig(cfg map[string]interface{}) error {
	skip, ok := cfg["skip_paths"]
	skipMap := map[string]struct{}{}

	if !ok {
		f.Skip = skipMap
		return errors.New("no 'skip_paths' key found")
	}
	sliceToMap(skip.([]interface{}), skipMap)
	f.Skip = skipMap

	return nil
}

func (f *FluentLoggerConfig) SetJWTClaimsConfig(cfg map[string]interface{}) error {
	claims, ok := cfg["include_jwt_claims"]
	claimsMap := map[string]struct{}{}

	if !ok {
		f.JWTClaims = claimsMap
		return errors.New("no 'include_jwt_claims' key found")
	}
	sliceToMap(claims.([]interface{}), claimsMap)
	f.JWTClaims = claimsMap

	return nil
}

func (f *FluentLoggerConfig) setBodyLoggingOptions(cfg map[string]interface{}) error {
	// 30 Mb
	defaultBodyLimit := int64(30)
	defaultAllowedContentTypes := map[string]struct{}{
		"application/json": {},
		"text/html":        {},
	}

	requestBodyLimit, err := getBodyLimit("request", cfg)
	if err != nil {
		f.Request.bodyLimit = defaultBodyLimit
	} else {
		f.Request.bodyLimit = int64(requestBodyLimit)
	}

	responseBodyLimit, err := getBodyLimit("response", cfg)
	if err != nil {
		f.Response.bodyLimit = defaultBodyLimit
	} else {
		f.Response.bodyLimit = int64(responseBodyLimit)
	}

	requestContentTypes, err := getAllowedContentType("request", cfg)
	if err != nil {
		f.Request.allowedContentTypes = defaultAllowedContentTypes
	} else {
		delete(requestContentTypes, "multipart/form-data")
		f.Request.allowedContentTypes = requestContentTypes
	}

	responseContentTypes, err := getAllowedContentType("response", cfg)
	if err != nil {
		f.Response.allowedContentTypes = defaultAllowedContentTypes
	} else {
		delete(requestContentTypes, "multipart/form-data")
		f.Response.allowedContentTypes = responseContentTypes
	}

	return nil
}

func getBodyLimit(key string, cfg map[string]interface{}) (int, error) {
	keyConfig, ok := cfg[key]
	if !ok {
		return 0, errors.New(fmt.Sprintf("no '%s' key found. using default fluent response config", key))
	}
	keyConfigMap, ok := keyConfig.(map[string]interface{})
	if !ok {
		return 0, errors.New(
			fmt.Sprintf("can't convert '%s' config to right type. using default fluent config", key),
		)
	}

	bodyLimitKey := "body_limit"
	_, ok = keyConfigMap[bodyLimitKey]
	if !ok {
		return 0, errors.New(fmt.Sprintf("response %v uses default values", bodyLimitKey))
	}

	return ConvertToInt(bodyLimitKey, keyConfigMap), nil
}

func getAllowedContentType(key string, cfg map[string]interface{}) (map[string]struct{}, error) {
	keyConfig, ok := cfg[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no '%s' key found. using default fluent response config", key))
	}
	keyConfigMap, ok := keyConfig.(map[string]interface{})
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("can't convert '%s' config to right type. using default fluent config", key),
		)
	}

	contentTypesKey := "allowed_content_types"
	contentTypes, ok := keyConfigMap[contentTypesKey]
	contentTypesMap := map[string]struct{}{}
	if !ok {
		return nil, errors.New(fmt.Sprintf("response %v uses default values", contentTypesKey))
	}
	sliceToMap(contentTypes.([]interface{}), contentTypesMap)

	return contentTypesMap, nil
}

func (f *FluentLoggerConfig) setMaskConfig(cfg map[string]interface{}) {
	key := "mask"

	maskConfig, ok := cfg[key]

	if !ok {
		printOutConfigError(key, errors.New(fmt.Sprintf("No masking config found")))
		return
	}

	f.Mask.Request = f.setMaskingConfig(maskConfig, "request")
	f.Mask.Response = f.setMaskingConfig(maskConfig, "response")
}

func (f *FluentLoggerConfig) setMaskingConfig(cfg interface{}, key string) map[string][]string {
	result := make(map[string][]string)
	maskConfigMap := cfg.(map[string]interface{})
	config, ok := maskConfigMap[key]

	if !ok {
		printOutConfigError(
			fmt.Sprintf("mask.%s", key), errors.New(fmt.Sprintf("No masking %s config found", key)),
		)
		return result
	}
	configConfigMap := config.(map[string]interface{})
	for k, v := range stringInterfaceToStringString(configConfigMap) {
		name := strings.Join([]string{key, k}, ".")
		result[name] = v
	}

	return result
}

func stringInterfaceToStringString(data map[string]interface{}) map[string][]string {
	requestMaskMap := map[string][]string{}
	for k1, v1 := range data {
		var maskSlice []string
		for _, v2 := range v1.([]interface{}) {
			maskSlice = append(maskSlice, v2.(string))
		}
		requestMaskMap[k1] = maskSlice
	}
	return requestMaskMap
}

func sliceToMap(skipSlice []interface{}, skipMap map[string]struct{}) {
	if length := len(skipSlice); length > 0 {
		for _, path := range skipSlice {
			skipMap[path.(string)] = struct{}{}
		}
	}
}
