package main

import (
	game "github.com/ntl-y/gameoflife/server"
	"github.com/ntl-y/gameoflife/server/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	hub := game.NewHub()
	go hub.Run()

	handler := handler.NewHandler(hub)

	srv := new(game.Server)
	if err := srv.Run(handler.InitRoutes()); err != nil {
		logrus.Fatal(err)
	}

}
