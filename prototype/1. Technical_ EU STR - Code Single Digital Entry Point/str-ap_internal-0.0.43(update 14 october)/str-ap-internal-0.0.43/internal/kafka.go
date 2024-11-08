package internal

import (
	"crypto/tls"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/spf13/viper"
	"time"
)

func getMechanism() plain.Mechanism {
	mechanism := plain.Mechanism{
		Username: viper.GetString("SASL_USERNAME"),
		Password: viper.GetString("SASL_PASSWORD"),
	}
	return mechanism
}

// Writer returns an authenticated kafka writer for a topic
func Writer(topic string) kafka.Writer {

	// Transports are responsible for managing connection pools and other resources,
	// it's generally best to create a few of these and share them across your
	// application.
	sharedTransport := &kafka.Transport{
		SASL: getMechanism(),
		TLS:  &tls.Config{},
	}

	w := kafka.Writer{
		Addr:                   kafka.TCP(viper.GetString("BOOTSTRAP_SERVERS")),
		Topic:                  topic,
		Balancer:               &kafka.RoundRobin{},
		Transport:              sharedTransport,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		Async:                  true,
		Logger:                 kafka.LoggerFunc(Infof),
		ErrorLogger:            kafka.LoggerFunc(Fatalf),
	}

	return w

}

// Reader returns an authenticated kafka reader for a topic & consumer-group
func Reader(topic string, groupID string) *kafka.Reader {

	dialer := &kafka.Dialer{
		SASLMechanism: getMechanism(),
		Timeout:       5 * time.Second,
		TLS:           &tls.Config{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{viper.GetString("BOOTSTRAP_SERVERS")},
		GroupID:        groupID,
		Topic:          topic,
		Dialer:         dialer,
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second, // flushes commits to Kafka every second
		Logger:         kafka.LoggerFunc(Infof),
		ErrorLogger:    kafka.LoggerFunc(Fatalf),
	})

	return r
}
