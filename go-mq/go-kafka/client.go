package gokafka

import (
	"fmt"
	"github.com/IBM/sarama"
	golog "github.com/gif-gif/go.io/go-log"
	"os"
	"strconv"
	"time"
)

type GoKafka struct {
	conf Config
	sarama.Client
}

func (cli *GoKafka) init() (err error) {
	id := strconv.Itoa(os.Getpid())
	config := sarama.NewConfig()
	config.ClientID = id
	if cli.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cli.conf.User
		config.Net.SASL.Password = cli.conf.Password
	}

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为Manual，指定分区发送消息
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	// 分区策略为Hash，解决相同key的消息落在一个分区
	//config.Producer.Partitioner = sarama.NewHashPartitioner
	// 分区策略为Random，解决消费组分布式部署
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true              // 自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	config.Consumer.Offsets.Retry.Max = 5
	if cli.conf.OffsetNewest {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
		sarama.NewBalanceStrategySticky(),
		sarama.NewBalanceStrategyRange(),
	}

	config.ChannelBufferSize = 1024
	//config.Version = sarama.V0_10_2_0

	config.Consumer.Group.Heartbeat.Interval = 5 * time.Second
	config.Consumer.Group.Session.Timeout = 15 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 12 * time.Second
	config.Producer.Timeout = 10 * time.Second

	if cli.conf.Timeout > 0 {
		config.Producer.Timeout = time.Duration(cli.conf.Timeout) * time.Second
	}

	if cli.conf.HeartbeatInterval > 0 {
		config.Consumer.Group.Heartbeat.Interval = time.Duration(cli.conf.HeartbeatInterval) * time.Second
	}
	if cli.conf.SessionTimeout > 0 {
		config.Consumer.Group.Session.Timeout = time.Duration(cli.conf.SessionTimeout) * time.Second
	}
	if cli.conf.RebalanceTimeout > 0 {
		config.Consumer.Group.Rebalance.Timeout = time.Duration(cli.conf.RebalanceTimeout) * time.Second
	}

	//config.Consumer.Group.InstanceId = id

	cli.Client, err = sarama.NewClient(cli.conf.Addrs, config)
	if err != nil {
		golog.WithTag("gokafka").Error(err)
	}

	return
}

func (cli *GoKafka) CreateTopicsRequest(topicName string, partitions int, replicationFactors int) error {
	request := &sarama.CreateTopicsRequest{}
	request.TopicDetails = make(map[string]*sarama.TopicDetail)
	request.TopicDetails[topicName] = &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicationFactors),
	}
	broker := cli.Brokers()[0]
	err := broker.Open(cli.Config())
	if err != nil {
		return err
	}
	defer broker.Close()
	ok, err := broker.Connected()
	if err != nil {
		return err
	}
	if ok {
		_, err = broker.CreateTopics(request)
		return err
	} else {
		return fmt.Errorf(" broker is not connected")
	}
}

func (cli *GoKafka) Close() {
	if !cli.Client.Closed() {
		cli.Client.Close()
	}
}

func (cli *GoKafka) Producer() iProducer {
	return &producer{GoKafka: cli, msg: &sarama.ProducerMessage{}}
}

func (cli *GoKafka) Consumer() iConsumer {
	return &consumer{GoKafka: cli}
}
