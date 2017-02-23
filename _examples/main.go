package main

import (
	"errors"

	"github.com/fundapps/logger"
)

func main() {
	infoOject := "Foo Info"
	logger.Info("Here is an info log", logger.Fields{"count": 1, "context": infoOject})

	err := errors.New("This is an error")
	logger.Error(err)

	logger.Error(logger.WrapErrorWithContext(err, "context error", logger.Fields{"code": 500}))

}
