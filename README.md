# krakend-fluentd-request-logger
Sends request logs to fluentd

#  Usage 


in router_engine.go of krakend-ce add FluentLoggerWithConfig middleware

```go
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
    if !cfg.Debug {
    gin.SetMode(gin.ReleaseMode)
    }

    engine := gin.New()

    engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: w}), gin.Recovery())
```

replace with

```go
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
			fluentd.FluentLoggerWithConfig(logger, cfg.ExtraConfig),
			gin.LoggerWithConfig(gin.LoggerConfig{Output: w}),
			gin.Recovery(),
		)
```