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

	// Запускаю в отдельной горутине, так как при deploy на heroku в течении 1 минуте питается установить порт
	// с переменной окружения, если не получается то приложения падает. А загрузка 1000 блоков занимает около 5 мин +-.
	go func(cfg *configs.Config) {
		if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
			logs.WithFields(logrus.Fields{
				"package":  "main",
				"function": "ListenAndServe",
				"error":    err,
			}).Fatal("The server is not up")
		}
	}(cfg)

	logs.Infof("Start server API on Port - %v", cfg.Port)

	// загрузка последних 1000 блоков
	lastLoadBlock, err := block.LoadBlocks()
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "main",
			"function": "LoadBlocks",
			"error":    err,
		}).Error("Error in func load block")
	}

	// вытаскиваем новые блоки, проверка 1 раз в секунду
	block.UpdateBlocks(lastLoadBlock)

}
