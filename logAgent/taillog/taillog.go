// 读取日志的模块

package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var (
	tailObj *tail.Tail
)

func Init(fileName string) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件那个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	tailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed :", err)
		return
	}
	fmt.Println("taillog init success")
	return
}

func ReadToChan() <-chan *tail.Line {
	return tailObj.Lines
}
