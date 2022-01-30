package api

var (
	BlockNumber   = "https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%v"
	BlockByNumber = "https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber" +
		"&tag=%v&boolean=true&apikey=%v"
)
