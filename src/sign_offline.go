package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bip39"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	mnemonicLength = 24

	algorithmSecp = 0
	algorithmEd   = 1

	secpString = "secp"
	edString   = "ed"
)

var (
	sigOutput  = os.Stderr
	infoOutput = os.Stdout
)

type args struct {
	mnemonic  string
	algorithm int
	noTrim    bool
}

func validateMnemonic(mnemonic string) {
	words := strings.Split(mnemonic, " ")
	if len(words) != mnemonicLength {
		panic(fmt.Sprintf("Invalid mnemonic length: Expected %d got %d", mnemonicLength, len(words)))
	}
}

func getSignBytes(noTrim bool) []byte {
	signBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil && err != io.EOF {
		panic("Failed to get sign bytes: " + err.Error())
	}

	if !noTrim {
		initialSize := len(signBytes)
		signBytes = []byte(strings.Trim(string(signBytes), "\n \r"))
		fmt.Fprintf(infoOutput, "%d bytes were trimmed", initialSize-len(signBytes))
	}

	fmt.Fprintf(infoOutput, "\nReady to sign %d bytes\n", len(signBytes))
	return signBytes
}

func getAlgorithm(alg string) int {
	switch alg {
	case secpString:
		return algorithmSecp
	case edString:
		return algorithmEd
	default:
		panic("Invalid private key algorithm")
	}
}

func getInputArgs() (args args) {
	mnemonicPtr := flag.String("mnemonic", "", "secret phrase")
	algPtr := flag.String("algorithm", "secp", "private key algorithm")
	noTrimPtr := flag.Bool("notrim", false, "don't trim sign bytes")
	flag.Parse()

	args.mnemonic = *mnemonicPtr
	args.mnemonic = strings.Trim(args.mnemonic, " ")
	validateMnemonic(args.mnemonic)

	args.algorithm = getAlgorithm(*algPtr)
	args.noTrim = *noTrimPtr
	return
}

func privateKeyFromMnemonic(mnemonic string, algorithm int) crypto.PrivKey {
	seed := bip39.MnemonicToSeed(mnemonic)
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, _ := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)

	switch algorithm {
	case algorithmSecp:
		return secp256k1.PrivKeySecp256k1(derivedPriv)
	case algorithmEd:
		return ed25519.GenPrivKeyFromSecret(derivedPriv[:])
	default:
		panic("Invalid private key algorithm")
	}
}

func sign(pkey crypto.PrivKey, signBytes []byte) []byte {
	signed, err := pkey.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	cdc := amino.NewCodec()
	sig, err := cdc.MarshalBinary(signed)
	if err != nil {
		panic(err)
	}

	return sig
}

func main() {
	inputArgs := getInputArgs()
	signBytes := getSignBytes(inputArgs.noTrim)

	pkey := privateKeyFromMnemonic(inputArgs.mnemonic, inputArgs.algorithm)

	signature := sign(pkey, signBytes)
	binary.Write(sigOutput, binary.BigEndian, signature)
}
