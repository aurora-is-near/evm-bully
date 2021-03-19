package nearapi

import (
	"encoding/hex"
	"testing"

	"github.com/near/borsh-go"
)

const signedTxHex = "1100000065766d2d62756c6c792e746573746e6574001a523c7d2d3434f82c0a069b883b83590fe23318ee491aeeba072b9b8e740713250000000000000014000000746573742d6163636f756e742e746573746e65747695cb1e847748021b0a7891410f458901bbbcf626e3eb6a27c1bf8e04d5e10f0100000003000040b2bac9e0191e02000000000000004ed8ac33b677616bed7430a21c190e71e1a80716cdfff9f972371ee0e08cc64c8b3fa1df79624ed987145a095be54e8a1ef38fb37640c397454619ea78f35605"

func TestSignedTransactionStruct(t *testing.T) {
	buf, err := hex.DecodeString(signedTxHex)
	if err != nil {
		t.Fatal(err)
	}
	var signedTx SignedTransaction
	if err = borsh.Deserialize(&signedTx, buf); err != nil {
		t.Fatal(err)
	}
}
