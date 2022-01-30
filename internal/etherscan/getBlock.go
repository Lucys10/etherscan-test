package etherscan

import (
	"apietherscan/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetBlockNumber(blockNumber string, apiKeyEther string) (int64, error) {
	url := fmt.Sprintf(blockNumber, apiKeyEther)

	r, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("Failed close body - ", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	var rb model.RespBlock
	if err := json.Unmarshal(data, &rb); err != nil {
		return 0, err
	}

	n, err := hexNumberToInt(rb.Result)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func GetBlockByNumber(blockByNumber string, hexValue string, apiKeyEther string) ([]byte, error) {
	url := fmt.Sprintf(blockByNumber, hexValue, apiKeyEther)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("Failed close body - ", err)
		}
	}()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
