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
		return "", fmt.Errorf("utils: cannot parse NEAR balance: %s", balance)
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

// formatWithCommas returns a human-readable value with commas.
func formatWithCommas(value string) string {
	res := ""
	for len(value) > 3 {
		res = "," + string(value[len(value)-3:]) + res
		value = value[:len(value)-3]
	}
	return value + res
}

// ParseNearAmount converts a human readable NEAR amount (potentially
// factional) to internal indivisible units. Effectively this multiplies given
// amount by NearNomination. Returns the parsed yoctoⓃ amount.
func ParseNearAmount(amount string) (string, error) {
	amount = cleanupAmount(amount)
	parts := strings.Split(amount, ".")
	wholePart := parts[0]
	fracPart := ""
	if len(parts) == 2 {
		fracPart = parts[1]
	}
	if len(parts) > 2 || len(fracPart) > NearNominationExp {
		return "", fmt.Errorf("utils: cannot parse as NEAR amount: %s", amount)
	}
	res := wholePart + fracPart + strings.Repeat("0", NearNominationExp-len(fracPart))
	res = strings.TrimLeft(res, "0")
	if res == "" {
		return "0", nil
	}
	return res, nil
}

// cleanupAmount removes commas from the input amount and returns the result.
func cleanupAmount(amount string) string {
	return strings.Replace(amount, ",", "", -1)
}

// ParseNearAmountAsBigInt converts a human readable NEAR amount (potentially
// factional) to internal indivisible units. Effectively this multiplies
// given amount by NearNomination. Returns the parsed yoctoⓃ amount.
func ParseNearAmountAsBigInt(amount string) (*big.Int, error) {
	res, err := ParseNearAmount(amount)
	if err != nil {
		return nil, err
	}
	var b big.Int
	_, ok := b.SetString(res, 10)
	if !ok {
		// Panic, because ParseNearAmount should never return a string that cannot
		// be converted to a big.Int.
		panic(fmt.Errorf("utils: cannot convert ParseNearAmount(%s) result to big.Int", amount))
	}
	return &b, nil
}
