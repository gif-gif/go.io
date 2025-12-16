package gokafka

import (
	"fmt"
	"os"
	"time"

	"github.com/IBM/sarama"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/samber/lo"
)

type GoKafka struct {
	conf Config
	sarama.Client
	redis                   *goredis.GoRedis
	ConsumerGroupInstanceId string
}

func (cli *GoKafka) GetConfig() Config {
	return cli.conf
}

func (cli *GoKafka) init() (err error) {
	createUniqueInstanceId := func() string {
		hostname, _ := os.Hostname()
		timestamp := time.Now().Unix()
		return fmt.Sprintf("consumer-%s-%d", hostname, timestamp)
	}

	//id := strconv.Itoa(os.Getpid())
	id := createUniqueInstanceId()

	cli.ConsumerGroupInstanceId = id
	config := sarama.NewConfig()
	config.ClientID = id
	config.ChannelBufferSize = cli.conf.ChannelBufferSize
	if cli.conf.Version == "" {
		config.Version = sarama.V3_6_0_0
	}
	if cli.conf.User != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cli.conf.User
		config.Net.SASL.Password = cli.conf.Password
	}

	if cli.conf.KeepAlive == 0 {
		cli.conf.KeepAlive = 10
	}

	config.Net.KeepAlive = time.Duration(cli.conf.KeepAlive) * time.Second

	// 等所有follower都成功后再返回
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 分区策略为Manual，指定分区发送消息
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	// 分区策略为Hash，解决相同key的消息落在一个分区
	//config.Producer.Partitioner = sarama.NewHashPartitioner
	// 分区策略为Random，解决消费组分布式部署
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = 1024 * 1024 * 100
	config.Producer.Return.Errors = true
	config.Consumer.Return.Errors = true

	//批量发送策略
	config.Producer.Flush.Messages = cli.conf.ProducerFlush.Messages                                     // 积累100条消息
	config.Producer.Flush.Bytes = cli.conf.ProducerFlush.Bytes                                           // 积累1MB数据
	config.Producer.Flush.Frequency = time.Duration(cli.conf.ProducerFlush.Frequency) * time.Millisecond // 每100ms刷新
	// 新增默认压缩方式
	config.Producer.Compression = sarama.CompressionLZ4

	// 增加是否自动提交的配置开启自动提交
	config.Consumer.Offsets.AutoCommit.Enable = cli.conf.AutoCommit.Enable                                  // 自动提交
	config.Consumer.Offsets.AutoCommit.Interval = time.Duration(cli.conf.AutoCommit.Interval) * time.Second // 间隔
	config.Consumer.Offsets.Retry.Max = 5
	if cli.conf.OffsetNewest {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategySticky(),
		sarama.NewBalanceStrategyRoundRobin(),
		sarama.NewBalanceStrategyRange(),
	}

	config.Consumer.Group.Heartbeat.Interval = time.Duration(cli.conf.ConsumerConfig.GroupConfig.HeartbeatInterval) * time.Second
	config.Consumer.Group.Session.Timeout = time.Duration(cli.conf.ConsumerConfig.GroupConfig.SessionTimeout) * time.Second
	config.Consumer.Group.Rebalance.Timeout = time.Duration(cli.conf.ConsumerConfig.GroupConfig.ReblanceInterval) * time.Second
	config.Consumer.Fetch.Default = cli.conf.ConsumerConfig.ConsumerFetchConfig.Default * 1024 * 1024
	config.Consumer.Fetch.Max = cli.conf.ConsumerConfig.ConsumerFetchConfig.Max * 1024 * 1024
	config.Consumer.Fetch.Min = cli.conf.ConsumerConfig.ConsumerFetchConfig.Min * 1024
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

	config.Consumer.Group.InstanceId = id

	cli.Client, err = sarama.NewClient(cli.conf.Addrs, config)
	if err != nil {
		golog.WithTag("gokafka").Error(err)
	}

	if cfg := cli.conf.RedisConfig; cfg.Addr != "" {
		cli.redis, err = goredis.New(cfg)
		if err != nil {
			golog.WithTag("gokafka").Error("Redis 初始化失败", err)
		}
	}

	return
}

//func (cli *GoKafka) CreateTopics(topics []string) (err error) {
//	kafkaTopics, err := Producer().Client().Topics()
//	if err != nil {
//		golog.WithTag("kafka-producer").Error(err)
//	}
//	for _, topic := range topics {
//		if lo.Contains(kafkaTopics, topic) {
//			continue
//		}
//		msg := KafkaMessageTest{
//			TopicName: topic,
//		}
//
//		_, _, err := gokafka.Producer().SendMessage(&msg)
//		if err != nil {
//			golog.Warn("KafkaMessageTest send message failed: ", err.Error())
//		}
//		fmt.Println("KafkaMessageTest send message success:" + topic)
//	}
//}

func (cli *GoKafka) CreateTopicRequest(topicName string, partitions int, replicationFactors int) error {
	kafkaTopics, err := Producer().Client().Topics()
	if err != nil {
		return err
	}

	if lo.Contains(kafkaTopics, topicName) {
		return nil
	}

	request := &sarama.CreateTopicsRequest{}
	request.TopicDetails = make(map[string]*sarama.TopicDetail)
	request.TopicDetails[topicName] = &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicationFactors),
	}
	broker := cli.Brokers()[0]
	err = broker.Open(cli.Config())
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

func (cli *GoKafka) CreateTopicsRequest(topicNames []string, partitions int, replicationFactors int) error {
	for _, topicName := range topicNames {
		err := cli.CreateTopicRequest(topicName, partitions, replicationFactors)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cli *GoKafka) Close() {
	if !cli.Client.Closed() {
		cli.Client.Close()
	}
}

// 消费者
func (cli *GoKafka) Consumer() IConsumer {
	return &consumer{GoKafka: cli}
}

// 生产者
func (cli *GoKafka) Producer(opts ...Option) IProducer {
	var focus bool
	for _, opt := range opts {
		switch opt.Name {
		case FocusName:
			focus = opt.Value.(bool)
		}
	}
	return &producer{GoKafka: cli, focus: focus}
}

func (c *GoKafka) GetKey(topic, msg string) string {
	return fmt.Sprintf("goio:mq:%s:%s", time.Now().Format("20060102"), goutils.MD5([]byte(topic+msg)))
}

func (c *GoKafka) Redis() *goredis.GoRedis {
	return c.redis
}

// 主题列表
func (c *GoKafka) Topics() []string {
	if c.Client == nil {
		return []string{}
	}

	topics, err := c.Client.Topics()
	if err != nil {
		golog.WithTag("gokafka").Error(err)
		return []string{}
	}

	return topics
}

// 分区数量
func (c *GoKafka) Partitions(topic string) []int32 {
	if c.Client == nil {
		return []int32{}
	}

	partitions, err := c.Client.Partitions(topic)
	if err != nil {
		golog.WithTag("gokafka").WithField("topic", topic).Error(err)
		return []int32{}
	}

	return partitions
}

// 分区数量
func (c *GoKafka) OffsetInfo(topic, groupId string) (data []map[string]int64) {
	data = []map[string]int64{}

	if c.Client == nil {
		return
	}

	partitions := c.Partitions(topic)
	if l := len(partitions); l == 0 {
		return
	}

	var (
		l = golog.WithTag("gokafka").WithField("groupId", groupId).WithField("topic", topic)
	)

	om, err := sarama.NewOffsetManagerFromClient(groupId, c.Client)
	if err != nil {
		l.Error(err)
		return
	}
	defer om.Close()

	for _, partition := range partitions {
		offset, err := c.Client.GetOffset(topic, partition, -1)
		if err != nil {
			l.Error(err)
			continue
		}

		pom, err := om.ManagePartition(topic, partition)
		if err != nil {
			l.Error(err)
			continue
		}

		nextOffset, msg := pom.NextOffset()
		if msg != "" {
			l.Error(msg)
			continue
		}

		backlog := offset
		if nextOffset != -1 {
			backlog -= nextOffset
		}

		data = append(data, map[string]int64{
			"partition":  int64(partition),
			"offset":     offset,
			"nextOffset": nextOffset,
			"backlog":    backlog,
		})
	}

	return
}
