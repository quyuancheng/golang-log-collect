package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 往kafka写日志的模块

var (
	client sarama.SyncProducer // 声明一个全局连接kafka的生产者client
)

// init初始化client

func Init(addrs []string) (err error) {
	config := sarama.NewConfig()
	// tailf包使用
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true                   //成功交付的消息将在success channel返回

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("kafka connect failed :", err)
		return
	}
	fmt.Println("kafka init success")
	return
}

func SendToKafka(topic, data string) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	//发送给kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send to kafka failed:", err)
		return
	}
	fmt.Printf("pid:%v, offset:%v\n", pid, offset)
	fmt.Println("发送kafka成功")
}
