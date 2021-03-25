package utils

import (
	"fmt"
)

func getTxnId(res map[string]interface{}) string {
	tx, ok := res["transaction"].(map[string]interface{})
	if ok {
		hash, ok := tx["hash"].(string)
		if ok {
			return hash
		}
	}
	return ""
}

func PrettyPrintResponse(res map[string]interface{}) {
	txnId := getTxnId(res)
	if txnId != "" {
		fmt.Printf("Transaction Id %s\n", txnId)
		// TODO: print transaction URL (requires explorer URL from config)
		// printTransactionurl(txnId)
	}
}
