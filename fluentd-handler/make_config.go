package fluentd_krakend_handler

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent" //nolint:goimports
	"github.com/luraproject/lura/logging"
	"strconv"
	"time"
)

type FluentLoggerConfig struct {
	FluentTag    string
	FluentConfig fluent.Config
	Skip         []string
	logger       logging.Logger
}

func ConvertToString(key string, cfg map[string]interface{}) string {
	value, ok := cfg[key]
	if ok {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

func ConvertToInt(key string, cfg map[string]interface{}) int {
	value, ok := cfg[key]
	valueString := fmt.Sprintf("%v", value)
	if ok {
		p, err := strconv.Atoi(valueString)
		if err == nil {
			return p
		}
	}
	return 0
}

func ConvertToBool(key string, cfg map[string]interface{}) bool {
	value, ok := cfg[key]
	if ok {
		p, ok := value.(bool)
		if ok {
			return p
		}
	}
	return false
}

func (f *FluentLoggerConfig) SetFluentHost(cfg map[string]interface{}) {
	f.FluentConfig.FluentHost = ConvertToString("fluent_host", cfg)
}

func (f *FluentLoggerConfig) SetFluentPort(cfg map[string]interface{}) {
	f.FluentConfig.FluentPort = ConvertToInt("fluent_port", cfg)
}

func (f *FluentLoggerConfig) SetFluentNetwork(cfg map[string]interface{}) {
	f.FluentConfig.FluentNetwork = ConvertToString("fluent_network", cfg)
}

func (f *FluentLoggerConfig) SetFluentSocketPath(cfg map[string]interface{}) {
	f.FluentConfig.FluentSocketPath = ConvertToString("fluent_socket_path", cfg)
}

func (f *FluentLoggerConfig) SetFluentTimeout(cfg map[string]interface{}) {
	f.FluentConfig.Timeout = time.Duration(ConvertToInt("timeout", cfg))
}

func (f *FluentLoggerConfig) SetFluentWriteTimeout(cfg map[string]interface{}) {
	f.FluentConfig.WriteTimeout = time.Duration(ConvertToInt("write_timeout", cfg))
}

func (f *FluentLoggerConfig) SetBufferLimit(cfg map[string]interface{}) {
	f.FluentConfig.BufferLimit = ConvertToInt("buffer_limit", cfg)
}

func (f *FluentLoggerConfig) SetRetryWait(cfg map[string]interface{}) {
	f.FluentConfig.RetryWait = ConvertToInt("retry_wait", cfg)
}

func (f *FluentLoggerConfig) SetMaxRetry(cfg map[string]interface{}) {
	f.FluentConfig.MaxRetry = ConvertToInt("max_retry", cfg)
}

func (f *FluentLoggerConfig) SetMaxRetryWait(cfg map[string]interface{}) {
	f.FluentConfig.MaxRetryWait = ConvertToInt("max_retry_wait", cfg)
}

func (f *FluentLoggerConfig) SetTagPrefix(cfg map[string]interface{}) {
	f.FluentConfig.TagPrefix = ConvertToString("tag_prefix", cfg)
}

func (f *FluentLoggerConfig) SetAsync(cfg map[string]interface{}) {
	f.FluentConfig.Async = ConvertToBool("async", cfg)
}

func (f *FluentLoggerConfig) SetForceStopAsyncSend(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("force_stop_async_send", cfg)
}

func (f *FluentLoggerConfig) SetSubSecondPrecision(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("sub_second_precision", cfg)
}

func (f *FluentLoggerConfig) SetRequestAck(cfg map[string]interface{}) {
	f.FluentConfig.ForceStopAsyncSend = ConvertToBool("request_ack", cfg)
}

func (f *FluentLoggerConfig) SetTag(cfg map[string]interface{}) {
	f.FluentTag = ConvertToString("fluent_tag", cfg)
}

func (f *FluentLoggerConfig) SetFluentConfig(cfg map[string]interface{}) {
	fluentConfig, ok := cfg["fluent_config"]
	if !ok {
		f.logger.Error("no 'fluent_config' key found. using default fluent config")
		return
	}

	fluentConfigMap, ok := fluentConfig.(map[string]interface{})
	if !ok {
		f.logger.Error("can't convert config to right type. using default fluent config")
		return
	}

	f.SetFluentHost(fluentConfigMap)
	f.SetFluentPort(fluentConfigMap)
	f.SetFluentNetwork(fluentConfigMap)
	f.SetFluentSocketPath(fluentConfigMap)
	f.SetFluentTimeout(fluentConfigMap)
	f.SetFluentWriteTimeout(fluentConfigMap)
	f.SetBufferLimit(fluentConfigMap)
	f.SetRetryWait(fluentConfigMap)
	f.SetMaxRetry(fluentConfigMap)
	f.SetMaxRetryWait(fluentConfigMap)
	f.SetTagPrefix(fluentConfigMap)
	f.SetAsync(fluentConfigMap)
	f.SetForceStopAsyncSend(fluentConfigMap)
	f.SetSubSecondPrecision(fluentConfigMap)
	f.SetRequestAck(fluentConfigMap)
}

func (f *FluentLoggerConfig) SetSkipConfig(cfg map[string]interface{}) {
	skip, ok := cfg["skip"]
	var emptySlice []string

	if !ok {
		f.logger.Debug("no 'skip' key found.")
		f.Skip = emptySlice
	}

	var skipSlice []string
	for _, param := range skip.([]interface{}) {
		skipSlice = append(skipSlice, param.(string))
	}
	f.Skip = skipSlice

	f.logger.Debug("krakend-fluentd-request-logger: 'skip' paths set")
}
