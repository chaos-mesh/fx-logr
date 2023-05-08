# fx-logr

Fx fxevent logr adapter. See <https://pkg.go.dev/go.uber.org/fx#WithLogger> for more details.

## Installation

```bash
go get github.com/chaos-mesh/fx-logr
```

## How to use

```go
import (
  "go.uber.org/fx"

  fxlogr "github.com/chaos-mesh/fx-logr"
)

func main() {
  fx.New(
    fx.WithLogger(
      fxlogr.WithLogr(&logger)
    ),
```

## License

Licensed under the Apache License, Version 2.0.

This adapter is inspired by the fxevent zap adapter (<https://github.com/uber-go/fx/blob/v1.19.2/fxevent/zap.go>) and <https://github.com/ipfans/fxlogger> (zerolog adapter). Thanks for the great work!
