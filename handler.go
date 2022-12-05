package handler

import (
	"errors"
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/v2/config"
)

const Namespace = "github_com/dmitrykaramin/krakend-fluentd-request-logger"

func EmptyFunc(_ *gin.Context) {}

func ReadConfig(conf *FluentLoggerConfig, extra config.ExtraConfig) error {
	appConfig, ok := extra[Namespace]

	if !ok {
		return errors.New("no app config found")
	}

	appConfigMap, ok := appConfig.(map[string]interface{})
	if !ok {
		return errors.New("can't convert config to right type")
	}

	err := conf.SetFluentConfig(appConfigMap)
	if err != nil {
		printOutError("fluentd config", err, "set %s error: %v \n")
	}
	err = conf.SetSkipConfig(appConfigMap)
	if err != nil {
		printOutError("fluentd 'skip' paths", err, "set %s error: %v \n")
	}
	err = conf.SetJWTClaimsConfig(appConfigMap)
	if err != nil {
		printOutError("fluentd 'include_jwt_claims' ", err, "set %s error: %v \n")
	}
	err = conf.setBodyLoggingOptions(appConfigMap)
	if err != nil {
		printOutError("fluentd 'response' ", err, "set %s error: %v \n")
	}

	conf.setMaskConfig(appConfigMap)

	return nil
}

func FluentLoggerWithConfig(logger logging.Logger, cfg config.ExtraConfig) gin.HandlerFunc {
	conf := FluentLoggerConfig{logger: logger}

	err := ReadConfig(&conf, cfg)
	if err != nil {
		logger.Error("krakend-fluentd-request-logger: ", err.Error())
		return EmptyFunc
	}

	return func(c *gin.Context) {
		fluentLogger, err := fluent.New(conf.FluentConfig)
		if err != nil {
			logger.Error(err)
			return
		}

		logWriter, err := NewLogWriter(c)
		if err != nil {
			logger.Error(err)
			return
		}
		c.Request.Header.Set("X-Correlation-ID", fmt.Sprint(uuid.New()))
		path := c.Request.URL.Path

		c.Next()

		if _, ok := conf.Skip[path]; ok {
			return
		}

		logWriter.CompleteLogData(c)
		data := logWriter.MakeLogData(conf)
		err = AddJwtData(data, conf.JWTClaims, c.Request.Header.Get("Authorization"))
		if err != nil {
			logger.Debug(err)
			return
		}

		err = fluentLogger.Post(conf.FluentTag, data)
		if err != nil {
			logger.Critical(err)
			return
		}

		err = fluentLogger.Close()
		if err != nil {
			logger.Error(err)
			return
		}
	}
}
