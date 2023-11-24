package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"os"
	"wbl0/internal"
)

func main() {
	cfg := internal.MustConfig()
	nc, err := nats.Connect(cfg.NatsConfig.Url)
	if err != nil {
		panic(fmt.Errorf("Error connecting to nats in pub: %s", err))
	}
	defer nc.Close()
	sc, err := stan.Connect(cfg.NatsConfig.StanClusterID, cfg.NatsConfig.ClientID, stan.NatsConn(nc))
	if err != nil {
		panic(fmt.Errorf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s",
			err, cfg.NatsConfig.Url))
	}
	defer sc.Close()
	data, err := os.ReadFile("model.json")
	if err != nil {
		panic(fmt.Errorf("Error reading file: %s", err))
	}
	err = sc.Publish("orders", data)
	if err != nil {
		panic(fmt.Errorf("Error publishing message: %s", err))
	}
	fmt.Println("Message sent")
}
