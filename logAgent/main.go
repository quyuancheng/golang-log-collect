package main

import (
	"11111/Log-collection/logAgent/config"
	"11111/Log-collection/logAgent/etcd"
	"11111/Log-collection/logAgent/kafka"
	"11111/Log-collection/logAgent/taillog"
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

var (
	cfg = new(config.AppConf)
)

func main() {
	// -----加载配置文件-----
	err := ini.MapTo(cfg, "./config/config.ini")
	if err != nil {
		fmt.Println("Load Config failed :", err)
	}

	// -----初始化kafka-----
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.ChanMaxSize)
	if err != nil {
		fmt.Println("init kafka failed :", err)
	}

	// -----初始化etcd-----
	err = etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd failed :", err)
	}
	// put key value to etcd test
	value := `[{"path":"F:\\GOProject\\src\\11111\\Log-collection\\logAgent\\game.log","topic":"game_log"}]`
	etcd.PutConf("qyc", value)
	// 1、-----从etcd获取日志收集项的配置信息-----
	logConfEntries, err := etcd.GetConf("qyc")
	if err != nil {
		fmt.Println("etcd GetConf failed :", err)
		return
	}
	fmt.Println("get logConf from etcd success :", logConfEntries)
	// 遍历logConfEntries，实例化每一个tailTask
	taillog.Init(logConfEntries)

	// 2、-----派一个哨兵去监视日志收集项的变化（有变化即时通知我的logAgent实现热加载配置）-----
	time.Sleep(time.Second * 100)
	// 运行
	//run()

}
