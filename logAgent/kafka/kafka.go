package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

// 往kafka写日志的模块

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer // 声明一个全局连接kafka的生产者client
	logDataChan chan *logData
)

// init初始化client

func Init(addrs []string, maxSize int) (err error) {
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

	//初始化logDataChan
	logDataChan = make(chan *logData, maxSize)
	// 启动一个发送kafka的任务
	go SendToKafka()
	return
}

// 把日志数据发送到一个chan中

func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg
}

func SendToKafka() {
	for {
		// 从chan中取出数据
		select {
		case dt := <-logDataChan:
			// 构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = dt.topic
			msg.Value = sarama.StringEncoder(dt.data)
			//发送给kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send to kafka failed:", err)
				return
			}
			fmt.Printf("pid:%v, offset:%v\n", pid, offset)
			fmt.Println("发送kafka成功")
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}
