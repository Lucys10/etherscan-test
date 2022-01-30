package etherscan

import (
	"apietherscan/api"
	"apietherscan/configs"
	"apietherscan/internal/model"
	"apietherscan/internal/store"
	"apietherscan/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	S             store.Store
	Logs          *logger.Log
	Cfg           *configs.Config
	QuantityBlock int64
}

func NewBlock(s store.Store, logs *logger.Log, cfg *configs.Config, quantityBlock int64) *Block {
	return &Block{
		S:             s,
		Logs:          logs,
		Cfg:           cfg,
		QuantityBlock: quantityBlock,
	}
}

func (b *Block) LoadBlocks() (int64, error) {

	// получаем текущий блок
	numBlock, err := GetBlockNumber(api.BlockNumber, b.Cfg.ApiKeyEther)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number - %v", err)
	}

	// цикл от текущего блока для загрузки 1000 последних блоков
	allTransInfo := make([]interface{}, 0)
	for i := numBlock - 1; i >= numBlock-b.QuantityBlock; i-- {

		hexValue := fmt.Sprintf("0x%x", i)

		// получаем список транзакций по номеру блока
		data, err := GetBlockByNumber(api.BlockByNumber, hexValue, b.Cfg.ApiKeyEther)
		if err != nil {
			return 0, fmt.Errorf("failed to get block by number - %v", err)
		}

		var rbn model.RespBlockByNumber
		if err := json.Unmarshal(data, &rbn); err != nil {
			return 0, fmt.Errorf("failed unmarshal - %v", err)
		}

		// в цикле перебираем все транзакции и добавляем в слайс
		for _, v := range rbn.Result.Transactions {
			block, err := hexNumberToInt(v.BlockNumber)
			if err != nil {
				b.Logs.WithFields(logrus.Fields{
					"package":  "etherscan",
					"function": "LoadBlocks",
					"error":    err,
				}).Error("failed convert hex to number")
			}
			ti := model.TransInfo{
				IdTrans:  v.TransactionIndex,
				From:     v.From,
				To:       v.To,
				NumBlock: block,
				Value:    v.Value,
			}

			allTransInfo = append(allTransInfo, ti)
		}

		b.Logs.Info("load block - ", i)
		time.Sleep(time.Millisecond * 300)
	}

	if err := b.S.InsertTransInfo(allTransInfo); err != nil {
		return 0, fmt.Errorf("failed Insert TransInfo - %v", err)
	}

	b.Logs.Info("Successful last thousand load block")

	return numBlock, nil
}

func (b *Block) UpdateBlocks(lastLoadBlock int64) {

	// загрузка 1000 блоков занимает 5 мин +-, поэтому за это время появляются новые блоки и что бы их не
	// пропустить в отдельной горутине запускаем загрузку этих блоков
	go diffBetweenLoadUpdate(b.S, b.Logs, lastLoadBlock, b.Cfg)

	var checkBlockNum int64 = 0

	for {

		numBlock, err := GetBlockNumber(api.BlockNumber, b.Cfg.ApiKeyEther)
		if err != nil {
			b.Logs.WithFields(logrus.Fields{
				"package":  "etherscan",
				"function": "GetBlockNumber",
				"error":    err,
			}).Error("failed to get block number")
		}

		if checkBlockNum == 0 {
			checkBlockNum = numBlock
		}

		// проверка когда появиться новый блок, предыдущий загружается в базу
		if numBlock != checkBlockNum {
			hexValue := fmt.Sprintf("0x%x", checkBlockNum)

			data, err := GetBlockByNumber(api.BlockByNumber, hexValue, b.Cfg.ApiKeyEther)
			if err != nil {
				b.Logs.WithFields(logrus.Fields{
					"package":  "etherscan",
					"function": "GetBlockByNumber",
					"error":    err,
				}).Error("failed to get block by number")
			}

			var rbn model.RespBlockByNumber
			if err := json.Unmarshal(data, &rbn); err != nil {
				b.Logs.WithFields(logrus.Fields{
					"package":  "etherscan",
					"function": "Unmarshal",
					"error":    err,
				}).Error("failed Unmarshal")
			}

			allTransInfo := make([]interface{}, 0, len(rbn.Result.Transactions))
			for _, v := range rbn.Result.Transactions {
				block, err := hexNumberToInt(v.BlockNumber)
				if err != nil {
					log.Println(err)
				}
				ti := model.TransInfo{
					IdTrans:  v.TransactionIndex,
					From:     v.From,
					To:       v.To,
					NumBlock: block,
					Value:    v.Value,
				}
				allTransInfo = append(allTransInfo, ti)
			}

			if len(allTransInfo) != 0 {
				if err := b.S.InsertTransInfo(allTransInfo); err != nil {
					b.Logs.WithFields(logrus.Fields{
						"package":  "etherscan",
						"function": "InsertTransInfo",
						"error":    err,
					}).Error("failed insert to mongoDB")
				}
			}
			b.Logs.Infof("Successful load block number - %v", checkBlockNum)
			checkBlockNum = numBlock
		}

		b.Logs.Infof("Current block number - %v", numBlock)

		time.Sleep(time.Second)
	}
}

func diffBetweenLoadUpdate(s store.Store, logs *logger.Log, lastBlockNum int64, cfg *configs.Config) {

	n, err := GetBlockNumber(api.BlockNumber, cfg.ApiKeyEther)
	if err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "etherscan",
			"function": "GetBlockByNumber",
			"error":    err,
		}).Error("failed to get block number")
	}

	allTransInfo := make([]interface{}, 0)
	for i := lastBlockNum; i < n; i++ {
		fmt.Println("Diff block - ", i)
		hexValue := fmt.Sprintf("0x%x", i)

		data, err := GetBlockByNumber(api.BlockByNumber, hexValue, cfg.ApiKeyEther)
		if err != nil {
			logs.WithFields(logrus.Fields{
				"package":  "etherscan",
				"function": "GetBlockByNumber",
				"error":    err,
			}).Error("failed to get block by number")
		}

		var rbn model.RespBlockByNumber
		if err := json.Unmarshal(data, &rbn); err != nil {
			logs.WithFields(logrus.Fields{
				"package":  "etherscan",
				"function": "Unmarshal",
				"error":    err,
			}).Error("failed Unmarshal")
		}

		for _, v := range rbn.Result.Transactions {
			block, err := hexNumberToInt(v.BlockNumber)
			if err != nil {
				logs.WithFields(logrus.Fields{
					"package":  "etherscan",
					"function": "LoadBlocks",
					"error":    err,
				}).Error("failed convert hex to number")
			}
			ti := model.TransInfo{
				IdTrans:  v.TransactionIndex,
				From:     v.From,
				To:       v.To,
				NumBlock: block,
				Value:    v.Value,
			}

			allTransInfo = append(allTransInfo, ti)
		}

		time.Sleep(time.Millisecond * 300)

	}

	if err := s.InsertTransInfo(allTransInfo); err != nil {
		logs.WithFields(logrus.Fields{
			"package":  "etherscan",
			"function": "InsertTransInfo",
			"error":    err,
		}).Error("failed insert to mongoDB")
	}
	logs.Info("Successful different block number")
}

func hexNumberToInt(hexStr string) (int64, error) {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	num, err := strconv.ParseInt(numStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}
