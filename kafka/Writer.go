package kafka

import (
	"log"
	"strings"
	"time"

	kafka "github.com/Shopify/sarama"

	"gitlab.zalopay.vn/bankintergration/offline-funding/offline-funding-common/logger"
)

//Writer ...
type Writer interface {
	WriteRaw([]byte)
	Write(kafka.Encoder)
}

type writer struct {
	topic    string
	producer kafka.AsyncProducer
}

//CreateWriter ....
func CreateWriter(addrStr, topic string) Writer {
	cfg := kafka.NewConfig()
	cfg.Producer.RequiredAcks = kafka.WaitForAll
	cfg.Producer.Flush.Frequency = 50 * time.Millisecond

	//addrs := strings.Split("10.205.21.47:9092", ",")
	addrs := strings.Split(addrStr, ",")
	logger.Infof("CreateWriter: %v", addrs)
	producer, err := kafka.NewAsyncProducer(addrs, cfg)
	if err != nil {
		log.Println(err)
	}
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write entry:", err)
		}
	}()
	return &writer{topic: topic, producer: producer}
}
func (w *writer) Write(v kafka.Encoder) {
	w.producer.Input() <- &kafka.ProducerMessage{
		Topic: w.topic,
		Value: v,
	}
}
func (w *writer) WriteRaw(v []byte) {
	w.producer.Input() <- &kafka.ProducerMessage{
		Topic: w.topic,
		Value: kafka.ByteEncoder(v),
	}
}
