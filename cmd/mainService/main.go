package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	internalhttp "github.com/PalPalych7/OtusProjectWork/internal/server/http"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	var logg ms.Logger
	var myBandid ms.MyBandit
	var storage ms.Storage
	var server ms.Server
	var err error

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	flag.Parse()
	config := NewConfig(configFile)

	logg = logger.New(config.Logger.LogFile, config.Logger.Level)
	logg.Info("Start!")
	myBandid = manyarmedbandit.New(config.Bandit)

	storage = sqlstorage.New(config.DB, myBandid)
	ctxDB, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.HTTP.TimeOutSec))
	defer cancel()
	if err = storage.Connect(ctxDB); err != nil {
		logg.Error(err.Error())
		time.Sleep(time.Minute * 1)
	} else {
		logg.Info("successful connect to DB")
	}
	defer storage.Close()

	server = internalhttp.NewServer(ctx, storage, config.HTTP, logg)
	defer server.Stop()

	go func() {
		if err := server.Serve(); err != nil {
			logg.Fatal("failed to start http server: " + err.Error())
		} else {
			logg.Info("Server was started")
		}
		<-ctx.Done()
	}()
	<-ctx.Done()
	logg.Info("Finish!")
}
