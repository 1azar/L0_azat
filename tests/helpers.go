package tests

import (
	"L0_azat/internal/domain"
	"encoding/json"
	"github.com/brianvoe/gofakeit"
	"github.com/nats-io/nats.go"
	"log"
	"math/rand"
)

func GenerateMsg() domain.Message {
	addr := gofakeit.Address()

	itemCount := rand.Intn(9) + 1
	items := make([]domain.Item, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = domain.Item{
			ChrtId:      rand.Intn(8999999) + 1000000,
			TrackNumber: "WBILMTESTTRACK",
			Price:       rand.Intn(100000),
			Rid:         gofakeit.UUID(),
			Name:        gofakeit.Word(),
			Sale:        rand.Intn(99) + 1,
			Size:        "0",
			TotalPrice:  rand.Intn(100000),
			NmId:        rand.Intn(8999999) + 1000000,
			Brand:       gofakeit.LastName(),
			Status:      202,
		}
	}

	var testMsg domain.Message = domain.Message{
		OrderUid:    gofakeit.UUID(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		DeliveryInfo: domain.DeliveryInfo{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: addr.Street + "15",
			Region:  addr.State,
			Email:   gofakeit.Email(),
		},
		PaymentInfo: domain.PaymentInfo{
			Transaction:  gofakeit.UUID(),
			RequestId:    "",
			Currency:     gofakeit.CurrencyShort(),
			Provider:     "wbpay",
			Amount:       99,
			PaymentDt:    1637907727,
			Bank:         "1637907727",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        gofakeit.UUID(),
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       gofakeit.Date(),
		OofShard:          "1",
	}

	return testMsg
}

func SendMsg(url, subject string, msg domain.Message) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	jsonMsg, err := json.Marshal(msg)

	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	err = nc.Publish(subject, jsonMsg)
	if err != nil {
		log.Fatal("Error publishing message:", err)
	}
}
