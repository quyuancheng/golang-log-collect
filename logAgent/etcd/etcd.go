package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// 初始化Init
var (
	cli *clientv3.Client
)

// 日志条目

type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

func Init(addr string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connect to etcd failed :", err)
		return
	}
	fmt.Println("Connect to etcd success")
	return
}

// 从etcd中根据Key获取配置项

func GetConf(key string) (logConfEntry []*LogEntry, err error) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed :", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err := json.Unmarshal(ev.Value, &logConfEntry)
		if err != nil {
			fmt.Println("json unmarshal failed :", err)
		}
	}
	return
}

// 在etcd中设置key:value (test API)

func PutConf(key string, val string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		fmt.Println("put key val error")
		return
	}
	fmt.Println("put to etcd success")
}
