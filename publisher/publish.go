package main

import (
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"math/rand"
	"wbl0/internal"
	"wbl0/internal/model"
)

func main() {
	cfg := internal.MustConfig()
	nc, err := nats.Connect(cfg.NatsConfig.Url)
	if err != nil {
		panic(fmt.Errorf("Error connecting to nats in pub: %s", err))
	}
	defer nc.Close()
	sc, err := stan.Connect(cfg.NatsConfig.StanClusterID, "publisher", stan.NatsConn(nc))
	if err != nil {
		panic(fmt.Errorf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s",
			err, cfg.NatsConfig.Url))
	}
	defer sc.Close()

	var data []model.Order
	for i := 0; i < 1000; i++ {
		data = append(data, model.Order{
			OrderUid:    fmt.Sprintf("testUid%d", i),
			TrackNumber: fmt.Sprintf("testTrackNumber-%d", i),
			Entry:       fmt.Sprintf("testEntry-%d", i),
			Delivery: model.Delivery{
				Name:    gofakeit.Name(),
				Phone:   gofakeit.Phone(),
				Zip:     gofakeit.Zip(),
				City:    gofakeit.City(),
				Address: gofakeit.Street(),
				Region:  gofakeit.State(),
				Email:   gofakeit.Email(),
			},
			Payment: model.Payment{
				Transaction:  fmt.Sprintf("testTransaction-%d", i),
				RequestID:    fmt.Sprintf("testPaymentRequestID-%d", i),
				Currency:     "RUB",
				Provider:     gofakeit.Word(),
				Amount:       rand.Intn(1000),
				PaymentDt:    rand.Intn(1000),
				Bank:         gofakeit.Word(),
				DeliveryCost: rand.Intn(1000),
				GoodsTotal:   rand.Intn(1000),
				CustomFee:    rand.Intn(1000),
			},
			Items:             testItems(),
			Locale:            gofakeit.Word(),
			InternalSignature: gofakeit.Word(),
			CustomerId:        fmt.Sprintf("testCustomerID-%d", i),
			DeliveryService:   gofakeit.Word(),
			Shardkey:          gofakeit.Word(),
			SmId:              rand.Intn(1000),
			DateCreated:       gofakeit.ConnectiveTime(),
			OofShard:          "testOofShard",
		})
	}
	for _, v := range data {
		testOrder, err := json.Marshal(v)
		if err != nil {
			fmt.Errorf("Error marshalling order: %s", err)
		}
		err = sc.Publish("orders", testOrder)
		if err != nil {
			panic(fmt.Errorf("Error publishing message: %s", err))
		}
	}
	//newExampleOrder, err := os.ReadFile("model.json")
	//if err != nil {
	//	fmt.Errorf("Error reading file: %s", err)
	//}
	//err = sc.Publish("orders", newExampleOrder)

}
func testItems() []model.Item {
	var testItems []model.Item
	for i := 0; i < rand.Intn(5); i++ {
		testItems = append(testItems, model.Item{
			ChrtID:      rand.Intn(1000),
			TrackNumber: gofakeit.Word(),
			Price:       rand.Intn(1000),
			Rid:         gofakeit.Word(),
			Name:        gofakeit.Word(),
			Sale:        rand.Intn(1000),
			Size:        gofakeit.Word(),
			TotalPrice:  rand.Intn(1000),
			NmID:        rand.Intn(1000),
			Brand:       gofakeit.Word(),
			Status:      rand.Intn(1000),
		})
	}
	return testItems
}
