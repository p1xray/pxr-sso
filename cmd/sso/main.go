package main

import "github.com/p1xray/pxr-sso/internal/app"

func main() {
	application := app.New()

	go application.Start()

	application.GracefulStop()
}
