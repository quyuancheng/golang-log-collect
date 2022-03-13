// 读取日志的模块

package taillog

import (
	"11111/Log-collection/logAgent/kafka"
	"fmt"
	"github.com/hpcloud/tail"
)

type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	// 根据路径打开日志文件
	tailObj.init()
	return
}

var err error

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件那个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed :", err)
		return
	}
	fmt.Println("taillog init success")
	go t.run() // 直接采集日志发送kafka
}

func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines: // 从tailObj的通道中一行一行读取日志
			// 发送kafka
			// kafka.SendToKafka(t.topic, line.Text)
			// 先把数据发送到一个通道中，kafka包中单独goroutine去通道中取数据发送kafka
			kafka.SendToChan(t.topic, line.Text)
			fmt.Println("send to kafka success", line.Text)
		}
	}
}
