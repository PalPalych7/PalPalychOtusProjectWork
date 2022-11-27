package main

import (
	"context"
	"encoding/json"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	rabbitmq "github.com/PalPalych7/OtusProjectWork/internal/rabbitMQ"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/statSenderConfig.toml", "Path to configuration file")
}

func main() {
	var logg ms.Logger
	var storage ms.Storage
	var myRQ ms.RabbitQueue
	var err error

	flag.Parse()
	config := NewConfig(configFile)
	logg = logger.New(config.Logger.LogFile, config.Logger.Level)
	logg.Info("Start!")
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	storage = sqlstorage.New(config.DB, nil)

	ctxDB, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Rabbit.TimeOutSec))
	defer cancel()

	if err = storage.Connect(ctxDB); err != nil {
		logg.Error(err.Error())
		time.Sleep(time.Second * time.Duration(config.Rabbit.SleepSecond))
	} else {
		logg.Info("successful connect to DB")
	}
	defer storage.Close()

	myRQ, err = rabbitmq.New(ctx, config.Rabbit)
	if err != nil {
		logg.Fatal("Error getting object for RQ", err.Error())
	}
	err = myRQ.Start()
	if err != nil {
		logg.Fatal(err.Error())
	}

	logg.Info("Connected to Rabit!")
	go func() {
		for {
			logg.Debug("I not sleep :).")
			// отправка оповещений
			myStatList, err2 := storage.GetBannerStat(ctx)
			countRec := len(myStatList)
			switch {
			case err2 != nil:
				logg.Error("Error in GetBannerStat", err2)
			case countRec == 0:
				logg.Debug("Nothing found for sending")
			default:
				logg.Debug("Found ", countRec, "record for sending")
				myMess, errMarsh := json.Marshal(myStatList)
				if errMarsh != nil {
					logg.Error("json.Marshal error", errMarsh)
				}
				if erSemdMess := myRQ.SendMess(myMess); erSemdMess != nil {
					logg.Error("Send mesage error", errMarsh)
				} else {
					logg.Info("message was succcessful send")
				}
				myStatID := myStatList[countRec-1].ID
				logg.Debug("max_stat_id=", myStatID)
				if errChID := storage.ChangeSendStatID(ctx, myStatID); errChID != nil {
					logg.Error("error in update max send ID -", errMarsh)
				}
			}
			time.Sleep(time.Second * time.Duration(config.Rabbit.SleepSecond))
		}
	}()
	<-ctx.Done()
	logg.Info("Finish")
}
