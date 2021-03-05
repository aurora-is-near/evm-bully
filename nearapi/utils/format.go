package utils

import (
	"fmt"
	"math/big"
	"strings"
)

// NearNominationExp defines the exponent for calculating how many indivisible
// units are there in one NEAR. See NearNomination.
const NearNominationExp = 24

// NearNomination defines the number of indivisible units in one NEAR. Derived
// from NearNominationExp.
var NearNomination = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(NearNominationExp), nil)

// FormatNearAmount converts account balance value from internal indivisible
// units to NEAR. 1 NEAR is defined by NearNomination. Effectively this
// divides given amount by NearNomination.
func FormatNearAmount(balance string) (string, error) {
	var bn big.Int
	_, suc := bn.SetString(balance, 10)
	if !suc {
		return "", fmt.Errorf("util: cannot parse NEAR balance: %s", balance)
	}
	balance = bn.String()
	wholeStr := "0"
	fractionStr := balance
	if len(balance) > NearNominationExp {
		wholeStr = balance[:len(balance)-NearNominationExp]
		fractionStr = string(balance[len(balance)-NearNominationExp:])
	} else {
		fractionStr = strings.Repeat("0", NearNominationExp) + fractionStr
		fractionStr = fractionStr[len(fractionStr)-NearNominationExp:]
	}
	res := formatWithCommas(wholeStr) + "." + fractionStr
	res = strings.TrimRight(res, "0")
	res = strings.TrimRight(res, ".")
	return res, nil
}

func formatWithCommas(value string) string {
	res := ""
	for len(value) > 3 {
		res = "," + string(value[len(value)-3:]) + res
		value = value[:len(value)-3]
	}
	return value + res
}
