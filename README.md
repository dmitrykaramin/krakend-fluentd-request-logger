# krakend-fluentd-request-logger

Sends request logs to fluentd

# Usage

#### in krakend.json

```json
{
  "version": 2,
  "extra_config": {
    "github_com/dmitrykaramin/krakend-fluentd-request-logger": {
      "fluent_config": {
        ...
      },
      "skip_paths": [
        "/path_to_skip"
      ],
      "include_jwt_claims": [
        "jwt_claim_field1",
        "jwt_claim_field2"
      ],
      "response": {
        "body_limit": 5,
        "allowed_content_types": [
          "application/json",
          "text/html"
        ]
      },
      "request": {
        "body_limit": 5,
        "allowed_content_types": [
          "application/json",
          "text/html"
        ]
      },
      "mask": {
        "request": {
          "headers": [
            "Authorization"
          ]
        },
        "response": {
          "body": ["access_token"]
        }
      }
    }
  },
  ...
}

```

## fluent_config

this json field accept following fluentd config fields:

`"fluent_port"`
`"fluent_host"`
`"fluent_network"`
`"fluent_socket_path"`
`"timeout"`
`"write_timeout"`
`"buffer_limit"`
`"retry_wait"`
`"max_retry"`
`"max_retry_wait"`
`"tag_prefix"`
`"async"`
`"force_stop_async_send"`
`"sub_second_precision"`
`"request_ack"`

absent fields will be filled out with default values:

`string` with `""`

`bool` with `false`

`int` with `0`

## skip_paths

is an array of strings: paths to skip from logging

## include_jwt_claims

is an array of jwt fields from jwt body to include in logging

## request/response

is an object of logging response options:

### body_limit
is symbols limit for logging - to prevent too large data logging. Default value - 5000

### allowed_content_types
is an array to define allowed content-type for logging. content-types not in array will not be logged.
Default value - ['application/json', 'html/text']

## mask
ability to mask sensitive data

### request/response
target to mask - can contain "body" or "headers";
"body" and "headers" can contain one level depth keys in body or header key to be masked

mask principle - if key's value contains more than 11 symbols they will be cut "first 4 ... last 4" 
if less than 11 symbols, value will be transformed in star "*" symbols


---

#### in router_engine.go of krakend-ce add FluentLoggerWithConfig middleware

```go
func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
    if !cfg.Debug {
    gin.SetMode(gin.ReleaseMode)
    }
    
    engine := gin.New()
    
    engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: w}), gin.Recovery())
    
    ...
}
```
or

```go
func NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	engine := luragin.NewEngine(cfg, opt)

	engine.NoRoute(opencensus.HandlerFunc(&config.EndpointConfig{Endpoint: "NoRoute"}, defaultHandler, nil))
	
	...
```

replace with

```go
import (
    ...
    "github.com/dmitrykaramin/krakend-fluentd-request-logger"
    ...
)


func NewEngine(cfg config.ServiceConfig, logger logging.Logger, w io.Writer) *gin.Engine {
    if !cfg.Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    
    engine := gin.New()
    engine.Use(
        handler.FluentLoggerWithConfig(logger, cfg.ExtraConfig),
        gin.LoggerWithConfig(gin.LoggerConfig{Output: w}),
        gin.Recovery(),
    )
    ...
}
```

or

```go
func NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
	engine := luragin.NewEngine(cfg, opt)
	engine.Use(handler.FluentLoggerWithConfig(opt.Logger, cfg.ExtraConfig))

	engine.NoRoute(opencensus.HandlerFunc(&config.EndpointConfig{Endpoint: "NoRoute"}, defaultHandler, nil))
```
