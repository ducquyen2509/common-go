package kafka

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	kafka "github.com/Shopify/sarama"

	"gitlab.zalopay.vn/bankintergration/offline-funding/offline-funding-common/logger"
)

//MultiReadCB ...
type MultiReadCB func(string, []byte) bool

//MultiReader use for consume multi topics
type MultiReader interface {
	ReadMulti(cb MultiReadCB)
}

type MulConsumerHandler struct {
	ready chan bool
	cb    MultiReadCB
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *MulConsumerHandler) Setup(kafka.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *MulConsumerHandler) Cleanup(kafka.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *MulConsumerHandler) ConsumeClaim(session kafka.ConsumerGroupSession, claim kafka.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		if consumer.cb(message.Topic, message.Value) {
			session.MarkMessage(message, "")
		}
	}

	return nil
}

//CreateMultiReader ...
func CreateMultiReader(addrStr, topics, group string, initOffset int64) MultiReader {
	config := kafka.NewConfig()
	config.Consumer.Offsets.Initial = initOffset
	config.Consumer.Return.Errors = true
	config.Version = kafka.V0_10_2_0
	// config.Group.Return.Notifications = true

	addrs := strings.Split(addrStr, ",")
	logger.Infof("CreateMultiReader, address: %v", addrs)

	consumerGroup, err := kafka.NewConsumerGroup(addrs, group, config)
	if err != nil {
		panic(err)
	}

	t := strings.Split(topics, ",")
	logger.Infof("CreateMultiReader, topics: %v", topics)

	return &reader{consumer: consumerGroup, topics: t}
}

func (r *reader) ReadMulti(cb MultiReadCB) {
	defer func() {
		if err := r.consumer.Close(); err != nil {
			logger.Errorf("closing client: %v", err)
		}
	}()
	logger.Infof("MultiReader topics: %v", r.topics)
	handler := MulConsumerHandler{
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

}
