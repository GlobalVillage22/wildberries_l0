package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"

	"time"
	"wbl0/internal"
	"wbl0/internal/db"
	"wbl0/internal/handler"
	"wbl0/internal/model"
	ch "wbl0/pkg/cash"
)

func manageMsg(m *stan.Msg, or *db.OrderRepo) {
	var order model.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		fmt.Errorf("Error unmarshalling order: %s", err)
		return
	}
	err = or.CreateOrder(order)
	if err != nil {
		fmt.Errorf("Error creating order: %s", err)
		return
	}
}
func printMsg(m *stan.Msg, i int) {
	fmt.Printf("[#%d] Received: %s\n", i, m)
}
func main() {
	cfg := internal.MustConfig()
	dbCon, err := db.NewConnecction(cfg)
	if err != nil {
		panic(fmt.Errorf("Error connecting to db: %s", err))
	}
	postgres := db.NewOrderRepo(dbCon)
	cash := ch.NewCash(10*time.Second, 100*time.Second)

	nc, err := nats.Connect(cfg.NatsConfig.Url)

	if err != nil {
		panic(fmt.Errorf("Error connecting to nats in sub: %s", err))
	}
	defer nc.Close()

	sc, err := stan.Connect(cfg.NatsConfig.StanClusterID, cfg.NatsConfig.ClientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			panic(fmt.Errorf("Connection lost, reason: %v", reason))
		}))
	if err != nil {
		panic(fmt.Errorf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s",
			err, cfg.NatsConfig.Url))
	}

	fmt.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", cfg.NatsConfig.Url,
		cfg.NatsConfig.StanClusterID, cfg.NatsConfig.ClientID)

	mcb := func(m *stan.Msg) {
		manageMsg(m, postgres)
	}
	//mcb2 := func(m *stan.Msg) {
	//	printMsg(m, 2)
	//}
	sub, err := sc.Subscribe(cfg.NatsConfig.Subject, mcb, stan.DurableName("wbl0"))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	fmt.Printf("Listening on [%s], clientID=[%s]\n", cfg.NatsConfig.Subject, cfg.NatsConfig.ClientID)
	//http.HandleFunc("/", handler.HandleOrder(postgres, cash))
	err = http.ListenAndServe("localhost:8080", handler.HandleOrder(postgres, cash))
	if err != nil {
		fmt.Errorf("Error listening: %s", err)
	}
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if cfg.NatsConfig.DurableName == "" {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
