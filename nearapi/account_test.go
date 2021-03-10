package nearapi

import (
	"path/filepath"
	"testing"
)

func TestReadAccessKey(t *testing.T) {
	var a Account
	receiverID := "test-account.testnet"
	fn := filepath.Join("testdata", receiverID+".json")
	err := a.readAccessKey(fn, receiverID)
	if err != nil {
		t.Fatal(err)
	}

	receiverID = "evm.test-account.testnet"
	fn = filepath.Join("testdata", receiverID+".json")
	err = a.readAccessKey(fn, receiverID)
	if err != nil {
		t.Fatal(err)
	}
}
