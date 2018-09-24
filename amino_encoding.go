package main

import (
	"bufio"
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

func validateMnemonic(mnemonic string) {
	words := strings.Split(mnemonic, " ")
	if len(words) != mnemonicLength {
		panic(fmt.Sprintf("Invalid mnemonic length: Expected %d got %d", mnemonicLength, len(words)))
	}
}

func getSignBytes() []byte {
	_, err := fmt.Fprintf(os.Stderr, "\nEnter bytes to sign:\n")

	inputReader := bufio.NewReader(os.Stdin)
	signBytes, err := inputReader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		panic("Failed to get sign bytes: " + err.Error())
	}
	fmt.Fprintf(os.Stderr, "\nReady to sign %d bytes\n", len(signBytes))
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

func processArgs() (mnemonic string, signBytes []byte, algorithm int) {
	mnemonicPtr := flag.String("mnemonic", "", "secret phrase")
	algPtr := flag.String("algorithm", "secp", "private key algorithm")
	flag.Parse()

	mnemonic = *mnemonicPtr
	mnemonic = strings.Trim(mnemonic, " ")
	fmt.Println(mnemonic)
	validateMnemonic(mnemonic)

	signBytes = getSignBytes()

	algorithm = getAlgorithm(*algPtr)
	return
}

func privateKeyFromMnemonic(mnemonic string, algorithm int) crypto.PrivKey {

	seed := bip39.MnemonicToSeed(mnemonic)
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, _ := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)

	return secp256k1.PrivKeySecp256k1(derivedPriv)

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
	mnemonic, signBytes, algorithm := processArgs()

	pkey := privateKeyFromMnemonic(mnemonic, algorithm)

	signature := sign(pkey, signBytes)
	ioutil.WriteFile("/Users/dmitry.kurilin/sig", signature, 0644)
	//	binary.Write(os.Stdout, binary.BigEndian, signature)
}
