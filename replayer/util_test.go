package replayer

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBigIntToRawU256(t *testing.T) {
	input := big.NewInt(256*256*256*13 + 37)
	expectedOutput := RawU256{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 0, 0, 37,
	}

	output, err := bigIntToRawU256(input)
	if err != nil {
		t.Fatalf("bigIntToRawU256(input) returns error, but it shouldn't")
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("bigIntToRawU256(input) result is %v, but expected value is %v", output, expectedOutput)
	}
}

func TestBigIntToRawU256Overflow(t *testing.T) {
	// 2^300 shouldn't fit into RawU256
	input := big.NewInt(0)
	input.Exp(big.NewInt(2), big.NewInt(300), nil)

	_, err := bigIntToRawU256(input)
	if err == nil {
		t.Fatalf("bigIntToRawU256(2^300) has to return error, but didn't")
	}
}
