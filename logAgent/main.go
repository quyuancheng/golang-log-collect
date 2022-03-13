package main

import (
	"11111/Log-collection/logAgent/config"
	"11111/Log-collection/logAgent/etcd"
	"11111/Log-collection/logAgent/kafka"
	"11111/Log-collection/logAgent/taillog"
	"11111/Log-collection/logAgent/utils"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var (
	cfg = new(config.AppConf)
	wg  sync.WaitGroup
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

	// 为了使每个logAgent都能拉去自己独有的配置，所以要以自己ip作为区分
	ip, err := utils.GetOutboundIP()
	if err != nil {
		fmt.Println("get ip failed:", err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ip)
	fmt.Println("etcdConfKey:", etcdConfKey)
	// 1、-----从etcd获取日志收集项的配置信息-----
	logConfEntries, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Println("etcd GetConf failed :", err)
		return
	}
	fmt.Println("get logConf from etcd success :", logConfEntries)

	// 遍历logConfEntries，实例化每一个tailTask,发送给Kafka
	taillog.Init(logConfEntries)

	// 2、-----派一个哨兵去监视日志收集项的变化（有变化即时通知我的logAgent实现热加载配置）-----
	newConChan := taillog.NewConfChan() // 从taillong中获取通道

	wg.Add(1)
	go etcd.WatchConf(cfg.EtcdConf.Key, newConChan) //哨兵发现最新的配置会通知上面的那个通道
	wg.Wait()
	// 运行
	//run()

}
