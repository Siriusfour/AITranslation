package destroy

import (
	"AITranslatio/Global"
	"os"
	"os/signal"
	"syscall"
)

func init() {

	go func() {

		c := make(chan os.Signal)

		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

		closeSignal := <-c

		Global.Logger.Info("关闭信号：", closeSignal.String())

		close(c)

		os.Exit(1)

	}()

}
