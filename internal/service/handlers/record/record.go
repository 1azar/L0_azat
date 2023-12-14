package record

import (
	"L0_azat/internal/domain"
	"encoding/json"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/nats-io/nats.go"
	"log/slog"
)

type msgSaver interface {
	SaveMsg(msg domain.Message) error
}

func New(log *slog.Logger, msgSaver msgSaver, cache *lru.Cache[string, any]) nats.MsgHandler {
	return func(msg *nats.Msg) {
		const fn = "handlers.record.record.New"
		log := log.With(
			slog.String("fn", fn))

		var receivedMsg domain.Message
		if err := json.Unmarshal(msg.Data, &receivedMsg); err != nil {
			log.Error("Error decoding json:", err)
			return
		}

		log.Debug(fmt.Sprintf("Received message: %v", receivedMsg))

		// db write
		if err := msgSaver.SaveMsg(receivedMsg); err != nil {
			log.Error("database write fail: ", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
		}

		// cache
		cache.Add(receivedMsg.OrderUid, receivedMsg)

		log.Debug("data has been written to the database")
	}
}
