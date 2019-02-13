package main

import (
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":9000"))
}

func hello(c echo.Context) error {
	output := make(chan bool, 1)
	errors := hystrix.Go("my_command", func() error {
		// talk to other services
		output <- true
		return nil
	}, nil)

	select {
	case <-output:
		return c.String(http.StatusOK, fmt.Sprintf("%s [%v]", "hello work!", output))
	case <-errors:
		return c.String(http.StatusBadGateway, fmt.Sprintf("%s [%v]", "error", errors))
	}

}
