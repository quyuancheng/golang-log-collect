package taillog

import (
	"11111/Log-collection/logAgent/etcd"
	"fmt"
	"time"
)

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntries  []*etcd.LogEntry     // 日志配置的列表
	tskMap      map[string]*TailTask //每一个任务
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntries []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntries:  logEntries,
		tskMap:      make(map[string]*TailTask, 16),
		newConfChan: make(chan []*etcd.LogEntry), // 无缓冲区的通道
	}
	for _, logEntry := range logEntries {
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		// 初始化的时候生成了多少tail task 都要记录下来，为了后续判断方便
		mk := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mk] = tailObj
	}
	go tskMgr.run()
}

// 监听自己的通道newConfChan，有了新的配置过来就做对应的处理

func (t *tailLogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("this is new conf: ", newConf)
			for _, conf := range newConf {
				mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mk]
				if ok {
					// 原来就有，不需要操作
					continue
				} else {
					// 新增的
					tailObj := NewTailTask(conf.Path, conf.Topic)
					mk := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
					t.tskMap[mk] = tailObj
				}
			}
			// 找出删除的
			for _, c1 := range t.logEntries {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					mk := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.tskMap[mk].cancelFunc()
				}
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

// 外部调用tskMgr的接口

func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
