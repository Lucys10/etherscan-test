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

	numBlock, err := GetBlockNumber(api.BlockNumber, b.Cfg.ApiKeyEther)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number - %v", err)
	}

	allTransInfo := make([]interface{}, 0)
	for i := numBlock - 1; i >= numBlock-b.QuantityBlock; i-- {

		hexValue := fmt.Sprintf("0x%x", i)

		data, err := GetBlockByNumber(api.BlockByNumber, hexValue, b.Cfg.ApiKeyEther)
		if err != nil {
			return 0, fmt.Errorf("failed to get block by number - %v", err)
		}

		var rbn model.RespBlockByNumber
		if err := json.Unmarshal(data, &rbn); err != nil {
			return 0, fmt.Errorf("failed unmarshal - %v", err)
		}

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

	b.Logs.Info("Successful load block")

	return numBlock, nil
}

func (b *Block) UpdateBlocks(lastLoadBlock int64) {

	go diffBetweenLoadUpdate(b.S, b.Logs, lastLoadBlock, b.Cfg)

	var checkBlockNum int64 = 0

	for {

		b.Logs.Info("Start update")
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

		b.Logs.Info("Chek block")
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

			fmt.Println("Load current block - ", checkBlockNum)
			fmt.Println("Len current block load allTransInfo - ", len(allTransInfo))
			checkBlockNum = numBlock
		}

		b.Logs.Infof("Successful load block number - %v", checkBlockNum)
		fmt.Println("Current block - ", numBlock)

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

	fmt.Println("Diff allTransInfo to Mongo - success!")
	fmt.Println("Len diff block allTransInfo - ", len(allTransInfo))
}

func hexNumberToInt(hexStr string) (int64, error) {
	numStr := strings.Replace(hexStr, "0x", "", -1)
	num, err := strconv.ParseInt(numStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}
