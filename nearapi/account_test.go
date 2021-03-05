package nearapi

import (
	"path/filepath"
	"testing"
)

func TestReadAccessKey(t *testing.T) {
	receiverID := "test-account.testnet"
	fn := filepath.Join("testdata", receiverID+".json")
	_, err := ReadAccessKey(fn, receiverID)
	if err != nil {
		t.Fatal(err)
	}

	receiverID = "evm.test-account.testnet"
	fn = filepath.Join("testdata", receiverID+".json")
	_, err = ReadAccessKey(fn, receiverID)
	if err != nil {
		t.Fatal(err)
	}
}
