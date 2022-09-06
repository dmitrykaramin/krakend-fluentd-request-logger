package handler

import (
	"errors"
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent" //nolint:goimports
	"github.com/luraproject/lura/logging"
	"strconv"
	"time"
)

type ResponseLoggerConfig struct {
	bodyLimit           int
	allowedContentTypes map[string]struct{}
}

type FluentLoggerConfig struct {
	FluentTag    string
	FluentConfig fluent.Config
	Skip         map[string]struct{}
	logger       logging.Logger
	JWTClaims    map[string]struct{}
	FlatHeaders  map[string]struct{}
	Response     ResponseLoggerConfig
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

func (f *FluentLoggerConfig) setResponseLoggingOptions(cfg map[string]interface{}) error {
	f.Response.bodyLimit = 5000
	f.Response.allowedContentTypes = map[string]struct{}{
		"application/json": {},
		"text/html":        {},
	}

	responseConfig, ok := cfg["response"]
	if !ok {
		return errors.New("no 'response' key found. using default fluent response config")
	}

	responseConfigMap, ok := responseConfig.(map[string]interface{})
	if !ok {
		return errors.New("can't convert response config to right type. using default fluent config")
	}

	f.setResponseBodyLimit(responseConfigMap)
	f.setResponseAllowedContentType(responseConfigMap)

	return nil
}

func (f *FluentLoggerConfig) setResponseBodyLimit(cfg map[string]interface{}) {
	key := "body_limit"

	_, ok := cfg[key]
	if !ok {
		printOutConfigError(key, errors.New(fmt.Sprintf("response %v uses default values", key)))
		return
	}

	f.Response.bodyLimit = ConvertToInt(key, cfg)
}

func (f *FluentLoggerConfig) setResponseAllowedContentType(cfg map[string]interface{}) {
	key := "allowed_content_types"

	contentTypes, ok := cfg[key]
	contentTypesMap := map[string]struct{}{}

	if !ok {
		printOutConfigError(key, errors.New(fmt.Sprintf("response %v uses default values", key)))
		return
	}

	sliceToMap(contentTypes.([]interface{}), contentTypesMap)
	f.Response.allowedContentTypes = contentTypesMap
}

func sliceToMap(skipSlice []interface{}, skipMap map[string]struct{}) {
	if length := len(skipSlice); length > 0 {
		for _, path := range skipSlice {
			skipMap[path.(string)] = struct{}{}
		}
	}
}
