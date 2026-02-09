package main

import (
	container "html2pdf/cmd/app"
)

func main() {
	app := container.InitApp()
	app.Init()
	app.Start()
}
