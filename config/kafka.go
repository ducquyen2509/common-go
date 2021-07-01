package config

type Kafka struct {
	Address, Group, Topic string
	Offset                int64
}
