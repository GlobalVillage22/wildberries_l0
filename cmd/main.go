package main

import (
	"fmt"
	_ "github.com/nats-io/nats.go"
	"net/http"
	"time"
	"wbl0/internal"
	"wbl0/internal/db"
	"wbl0/internal/handler"
	"wbl0/internal/nats"
	ch "wbl0/pkg/cash"
)

func main() {
	cfg := internal.MustConfig()
	postgres := db.NewOrderRepo(db.MustConnection(cfg))
	cash := ch.NewCash(10*time.Second, 100*time.Second)
	sc := nats.MustConnection(cfg)
	_, err := nats.NewSubscribe(cfg, sc, postgres)
	if err != nil {
		panic(fmt.Errorf("Error subscribing to nats: %s", err))
	}
	err = http.ListenAndServe("localhost:8080", handler.HandleOrder(postgres, cash))
	if err != nil {
		fmt.Errorf("Error listening: %s", err)
	}
}
