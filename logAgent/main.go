package main

import (
	"11111/Log-collection/logAgent/config"
	"11111/Log-collection/logAgent/kafka"
	"11111/Log-collection/logAgent/taillog"
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	cfg = new(config.AppConf)
)

func run() {
	for {
		select {
		case line := <-taillog.ReadToChan():
			// 发送到kafka
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
		}
	}
}

func main() {
	// 加载配置文件
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Println("Load Config failed :", err)
	}
	// 初始化kafka
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	// 初始化taillog
	err = taillog.Init(cfg.TaillogConf.Filename)
	// 运行
	run()
}
