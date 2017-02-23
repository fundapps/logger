FundApps Logger
================


Combined logger and Sentry handler.

## Install

```sh
  go get github.com/fundapps/logger
```

## Usage

Example
```go

package main

import (
	"errors"

	"github.com/fundapps/logger"
)

func main() {
	err := errors.New("This is an error")

	logger.Error(logger.WrapErrorWithContext(err, "context error", logger.Fields{"code": 500}))

}

```

## Development Usage

To disable Sentry during Development, use the environment variable `APP_ENV`
