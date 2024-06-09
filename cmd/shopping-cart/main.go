package main

import (
	"github.com/Pugpaprika21/pkg/shopping-cart/router"
	"github.com/Pugpaprika21/pkg/shopping-cart/server"
)

func main() {
	e := server.NewEchoServer()
	router.EchoRouter(e.Echo, e.Server)
	e.Echo.Logger.Fatal(e.Echo.Start(":" + e.Server.App.Port))
}
