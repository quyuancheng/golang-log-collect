package config

type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chan_max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"key"`
}

// ----- unused -----

type TaillogConf struct {
	Filename string `ini:"filename"`
}
