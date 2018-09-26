package main

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

type testCase struct {
	Mnemonic  string `json:"Mnemonic"`
	PubKey    string `json:"PubKey"`
	Msg       string `json:"Msg"`
	Algorithm string `json:"Algorithm"`
	Result    bool   `json:"Result"`
}

// type cases struct {
// 	Cases []testCase `json:"Cases"`
// }

var cases = []testCase{
	testCase{
		Mnemonic:  "prosper sing educate detail rose spoon cute few arctic focus recipe family sleep adapt sphere glue fantasy sorry limb rate pepper omit move limit",
		PubKey:    "cosmosaccpub1addwnpepqw2u4ke0duujy0ryyapv92g75m3pjjaw7dvv6jr0mv449e557dqpclyj642",
		Msg:       `{"account_number":"0","chain_id":"test-chain-mqSoIm","fee":{"amount":[{"amount":"0","denom":""}],"gas":"200000"},"memo":"","msgs":[{"inputs":[{"address":"cosmos1czgjqgtjrvlq37f3lqrph0tggzl5w7xfx5w9wc","coins":[{"amount":"100","denom":"mycoin"}]}],"outputs":[{"address":"cosmos1zukk5knm8gctmkzyytkwk4x58dge95y7c6lymg","coins":[{"amount":"100","denom":"mycoin"}]}]}],"sequence":"0"}`,
		Algorithm: "secp",
		Result:    true,
	},
	testCase{
		// different pubkeys
		Mnemonic:  "prosper sing educate detail rose spoon cute few arctic focus recipe family sleep adapt sphere glue fantasy sorry limb rate pepper omit move limit",
		PubKey:    "cosmosaccpub1addwnpepqv2k9ss7jg8kgwns4q9g8tjm309gugjxsp9j6yvgvl05xeuhjkquwxdr4w2",
		Msg:       `{"account_number":"0","chain_id":"test-chain-mqSoIm","fee":{"amount":[{"amount":"0","denom":""}],"gas":"200000"},"memo":"","msgs":[{"inputs":[{"address":"cosmos1czgjqgtjrvlq37f3lqrph0tggzl5w7xfx5w9wc","coins":[{"amount":"100","denom":"mycoin"}]}],"outputs":[{"address":"cosmos1zukk5knm8gctmkzyytkwk4x58dge95y7c6lymg","coins":[{"amount":"100","denom":"mycoin"}]}]}],"sequence":"0"}`,
		Algorithm: "secp",
		Result:    false,
	},
}

func verify(signature []byte, testCase testCase, t *testing.T) bool {
	cdc := amino.NewCodec()
	var unpackedSig []byte
	cdc.MustUnmarshalBinary(signature, &unpackedSig)
	pubKey, err := types.GetAccPubKeyBech32(testCase.PubKey)

	if err != nil {
		t.Errorf("Invalid PubKey string: %s", testCase.PubKey)
	}

	return pubKey.VerifyBytes([]byte(testCase.Msg), unpackedSig)
}

func TestCheck(t *testing.T) {
	for ind, testCase := range cases {
		signature := Sign(testCase.Mnemonic, []byte(testCase.Msg), testCase.Algorithm)
		result := verify(signature, testCase, t)
		if testCase.Result != result {
			t.Errorf("Testcase #%d failed. Result != %t", ind, testCase.Result)
		}
	}
}
