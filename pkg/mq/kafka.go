package mq

import (
	"github.com/IBM/sarama"
	"time"
)

type KafkaConfig struct {
	Scanner KafkaCluster `mapstructure:"scanner" json:"scanner" yaml:"scanner"`
	Web     KafkaCluster `mapstructure:"web" json:"web" yaml:"web"`
}

type KafkaCluster struct {
	BootstrapServers []string `mapstructure:"bootstrap-servers" json:"bootstrap-servers" yaml:"bootstrap-servers"`
	Consumer         Consumer `mapstructure:"consumer" json:"consumer" yaml:"consumer"`
}

type Consumer struct {
	GroupId          string `mapstructure:"group-id" json:"group-id" yaml:"group-id"`
	AutoOffsetReset  string `mapstructure:"auto-offset-reset" json:"auto-offset-reset" yaml:"auto-offset-reset"`
	EnableAutoCommit bool   `mapstructure:"enable-auto-commit" json:"enable-auto-commit" yaml:"enable-auto-commit"`
}

func NewKafkaManager(config *KafkaConfig) *KafkaClient {
	return &KafkaClient{
		Conf: config,
	}
}

type KafkaClient struct {
	Conf      *KafkaConfig
	kafkaName string
}

// 异步生产者 允许你将消息发送到 Kafka 集群，而不必等待确认。这对于高吞吐量场景非常有用。
func (client *KafkaClient) NewAsyncProducer() sarama.AsyncProducer {
	appConf := client.Conf.Web
	if client.kafkaName == "scanner" {
		appConf = client.Conf.Scanner
	}

	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Idempotent = true // 开启幂等性
	config.Net.MaxOpenRequests = 1    // 开启幂等性后 并发请求数也只能为1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = time.Duration(5) * time.Minute
	config.Producer.Transaction.ID = "txn_producer"
	config.Producer.MaxMessageBytes = 10 << 20 // 10 MB
	config.Consumer.Offsets.AutoCommit.Enable = appConf.Consumer.EnableAutoCommit
	// config.Consumer.Group.ResetInvalidOffsets = rawconfig.Web.Consumer.AutoOffsetReset
	config.Consumer.Group.InstanceId = appConf.Consumer.GroupId
	producer, _ := sarama.NewAsyncProducer(appConf.BootstrapServers, config)
	return producer
}

// 同步生产者 在发送消息时会等待 Kafka 确认，这对于需要保证消息确实发送成功的场景非常重要。
func (client *KafkaClient) NewSyncProducer() sarama.SyncProducer {
	appConf := client.Conf.Web
	if client.kafkaName == "scanner" {
		appConf = client.Conf.Scanner
	}

	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Idempotent = true // 开启幂等性
	config.Net.MaxOpenRequests = 1    // 开启幂等性后 并发请求数也只能为1
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = time.Duration(5) * time.Minute
	config.Producer.MaxMessageBytes = 10 << 20 // 10 MB
	// config.Consumer.Group.ResetInvalidOffsets = rawconfig.Web.Consumer.AutoOffsetReset
	producer, _ := sarama.NewSyncProducer(appConf.BootstrapServers, config)
	return producer
}

// 要链式调用
// name:"web" 或者 "scanner"
func (client *KafkaClient) KafkaName(name string) *KafkaClient {
	if name != "web" && name != "scanner" {
		panic("Kafka调用出错, name必须是web或scanner")
	}
	client.kafkaName = name
	return client
}
