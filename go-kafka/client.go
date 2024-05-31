package gokafka

import (
	"github.com/IBM/sarama"
	golog "github.com/gif-gif/go.io/go-log"
	"time"
)

type client struct {
	conf Config
	sarama.Client
}

func (cli *client) init() (err error) {
	config := sarama.NewConfig()

	if cli.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cli.conf.User
		config.Net.SASL.Password = cli.conf.Password
	}

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为Hash，解决相同key的消息落在一个分区
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true              // 自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	config.ChannelBufferSize = 1000
	config.Version = sarama.V0_10_2_0

	if cli.conf.Timeout > 0 {
		config.Producer.Timeout = time.Duration(cli.conf.Timeout) * time.Second
	}

	cli.Client, err = sarama.NewClient(cli.conf.Addrs, config)
	if err != nil {
		golog.WithTag("goo-kafka").Error(err)
	}

	return
}

func (cli *client) Close() {
	if !cli.Client.Closed() {
		cli.Client.Close()
	}
}
