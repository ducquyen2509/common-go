package kafka

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	kafka "github.com/Shopify/sarama"

	"github.com/ducquyen2509/common-go/logger"
)

//ReadCB ...
type ReadCB func([]byte) bool

//Reader ...
type Reader interface {
	Read(cb ReadCB)
}

type reader struct {
	consumer kafka.ConsumerGroup
	topics   []string
}

type ConsumerHandler struct {
	ready chan bool
	cb    ReadCB
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ConsumerHandler) Setup(cgs kafka.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerHandler) Cleanup(kafka.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *ConsumerHandler) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if consumer.cb(message.Value) {
			logger.Infof("Marked kafka message ( topic: %v, Partition: %v, Offset: %v )", message.Topic, message.Partition, message.Offset)
			session.MarkMessage(message, "")
		}
	}
	return nil
}

//CreateReader create single topic consumer
func CreateReader(addrStr, topic, group string, initOffset int64, partitions []int32) Reader {
	cfg := kafka.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Version = kafka.V0_10_2_0

	addrs := strings.Split(addrStr, ",")
	logger.Infof("CreateReader: %v", addrs)

	consumerGroup, err := kafka.NewConsumerGroup(addrs, group, cfg)
	if err != nil {
		panic(err)
	}

	return &reader{consumer: consumerGroup, topics: []string{topic}}
}

func (r *reader) Read(cb ReadCB) {
	handler := ConsumerHandler{
		ready: make(chan bool),
		cb:    cb,
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := r.consumer.Consume(ctx, r.topics, &handler); err != nil {
				logger.Errorf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			handler.ready = make(chan bool)
		}
	}()

	<-handler.ready // Await till the consumer has been set up

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info("terminating: context cancelled")
	case <-sig:
		logger.Info("terminating: via signal")
	}
	cancel()
	wg.Wait()

	if err := r.consumer.Close(); err != nil {
		logger.Errorf("Error closing client: %v", err)
		panic(err)
	}
}
