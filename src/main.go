package main

import (
	_ "bytes"
	_ "github.com/CodyGuo/godaemon"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"monitor-agent/config"
	"monitor-agent/kits"
	"monitor-agent/modules"
	"os"
	"time"
)

func main() {
	// 检查Agent进程是否重复启动
	if kits.CheckFile(config.PidFile) {
		f, err := ioutil.ReadFile(config.PidFile)
		if err != nil {
			kits.Log(err.Error(), "error", "main")
		}
		p, err := ps.FindProcess(int(kits.BytesToInt64(f)))
		if p != nil {
			println("monitor-agent进程已运行!")
			os.Exit(1)
		}
	}
	//默认睡眠时间60秒
	sleepInterval := 60
	for {
		pid := os.Getpid()
		err := ioutil.WriteFile(config.PidFile, kits.Int64ToBytes(int64(pid)), 0666)
		if err != nil {
			kits.Log(err.Error(), "error", "main")
		}
		//获取配置信息
		cnf, err := kits.GetConfig(config.Config{})
		if err != nil {
			// 配置文件异常进入睡眠模式
			sleepInterval += 60
			//最长睡眠时间30分钟
			if sleepInterval >= 1800 {
				sleepInterval = 60
			}
			time.Sleep(time.Duration(sleepInterval) * time.Second)
		} else {
			modules.Agent(cnf)
		}
	}
}
