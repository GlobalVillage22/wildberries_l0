package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"wbl0/internal"
	"wbl0/internal/db"
	"wbl0/internal/model"
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

func MustConnection(cfg *internal.Config) stan.Conn {
	nc, err := nats.Connect(cfg.NatsConfig.Url)

	if err != nil {
		panic(fmt.Errorf("Error connecting to nats in sub: %s", err))
	}

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
	return sc
}
func NewSubscribe(cfg *internal.Config, sc stan.Conn, or *db.OrderRepo) (stan.Subscription, error) {
	mcb := func(m *stan.Msg) {
		manageMsg(m, or)
	}
	sub, err := sc.Subscribe(cfg.NatsConfig.Subject, mcb, stan.DurableName(cfg.NatsConfig.DurableName))
	if err != nil {
		sc.Close()
		return nil, err
	}
	return sub, nil
}
