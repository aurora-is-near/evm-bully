package utils

import (
	"fmt"
)

func getTxnID(res map[string]interface{}) string {
	tx, ok := res["transaction"].(map[string]interface{})
	if ok {
		hash, ok := tx["hash"].(string)
		if ok {
			return hash
		}
	}
	return ""
}

// PrettyPrintResponse pretty prints some info from the given JSON result map res.
func PrettyPrintResponse(res map[string]interface{}) {
	txnID := getTxnID(res)
	if txnID != "" {
		fmt.Printf("Transaction Id %s\n", txnID)
		// TODO: print transaction URL (requires explorer URL from config)
		// printTransactionurl(txnId)
	}
}
