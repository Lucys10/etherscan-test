package main

import (
	"apietherscan/configs"
	"apietherscan/internal/etherscan"
	"apietherscan/internal/handlers"
	"apietherscan/internal/store"
	"apietherscan/pkg/db"
	"apietherscan/pkg/logger"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {

	logs := logger.NewLogger(logrus.InfoLevel)

	cfg, err := configs.GetConfig()
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "GetConfig",
			"error":    err,
		}).Fatal("failed get config")
	}

	fmt.Println(cfg)

	dbm, err := db.NewMongo(cfg.MongoPass)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "NewMongo",
			"error":    err,
		}).Fatal("failed connection to MongoDB")
	}

	s := store.NewStore(dbm)
	r := chi.NewRouter()
	h := handlers.Handlers{S: s}
	h.Register(r)

	block := etherscan.NewBlock(s, logs, cfg)

	lastLoadBlock, err := block.LoadBlocks()
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "LoadBlocks",
			"error":    err,
		}).Error("Error in func load block")
	}

	go block.UpdateBlocks(lastLoadBlock)

	logs.WithFields(logrus.Fields{
		"ServerAddress": ":8090",
		"Log_level":     "Info",
	}).Info("Start controller-service...")

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "ListenAndServe",
			"error":    err,
		}).Fatal("The server is not up")
	}
}
