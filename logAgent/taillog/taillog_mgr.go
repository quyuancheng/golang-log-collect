package taillog

import "11111/Log-collection/logAgent/etcd"

var tskMgr *tailLogMgr

type tailLogMgr struct {
	logEntries []*etcd.LogEntry     // 日志配置的列表
	tskMap     map[string]*TailTask //每一个任务
}

func Init(logEntries []*etcd.LogEntry) {
	tskMgr = &tailLogMgr{
		logEntries: logEntries,
	}
	for _, logEntry := range logEntries {
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
