package main

import (
	"apietherscan/configs"
	"apietherscan/internal/etherscan"
	"apietherscan/internal/handlers"
	"apietherscan/internal/store"
	"apietherscan/pkg/db"
	"apietherscan/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
)

const quantityBlock int64 = 1000

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
	h := handlers.Handlers{S: s, Logs: logs}
	h.Register(r)

	block := etherscan.NewBlock(s, logs, cfg, quantityBlock)

	go func() {
		if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
			logs.WithFields(logrus.Fields{
				"package":  "main",
				"function": "ListenAndServe",
				"error":    err,
			}).Fatal("The server is not up")
		}
	}()

	lastLoadBlock, err := block.LoadBlocks()
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "LoadBlocks",
			"error":    err,
		}).Error("Error in func load block")
	}

	block.UpdateBlocks(lastLoadBlock)

}
