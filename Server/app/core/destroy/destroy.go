package destroy

import (
	"AITranslatio/Global"
	"os"
	"os/signal"
	"syscall"
)

func init() {

	go func() {

		c := make(chan os.Signal) //定义一个用于接受系统信号的ch
		defer close(c)

		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		closeSignal := <-c //阻塞在这里，直到接收到ctrl+c

		Global.Logger.Info("关闭信号：", closeSignal.String())

		os.Exit(1) //退出系统

	}()

}
