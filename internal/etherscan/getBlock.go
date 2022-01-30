package etherscan

import (
	"apietherscan/internal/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetBlockNumber(blockNumber string, apiKeyEther string) (int64, error) {
	url := fmt.Sprintf(blockNumber, apiKeyEther)

	r, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	var rb model.RespBlock
	if err := json.Unmarshal(data, &rb); err != nil {
		log.Println(err)
	}

	n, err := hexNumberToInt(rb.Result)
	if err != nil {
		log.Println(err)
	}

	return n, nil
}

func GetBlockByNumber(blockByNumber string, hexValue string, apiKeyEther string) ([]byte, error) {
	url := fmt.Sprintf(blockByNumber, hexValue, apiKeyEther)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
