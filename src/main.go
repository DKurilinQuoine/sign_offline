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
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	mnemonicLength = 24

	secpString = "secp"
	edString   = "ed"
)

var (
	sigOutput  = os.Stderr
	infoOutput = os.Stdout
)

type args struct {
	mnemonic  string
	algorithm string
	noTrim    bool
}

// public

// Sign generates a private key based on mnemonic and algorithm and signs signBytes
func Sign(mnemonic string, signBytes []byte, algorithm string) []byte {
	pkey := privateKeyFromMnemonic(mnemonic, algorithm)
	return sign(pkey, signBytes)
}

func main() {
	inputArgs := getInputArgs()
	signBytes := getSignBytes(inputArgs.noTrim)
	signature := Sign(inputArgs.mnemonic, signBytes, inputArgs.algorithm)
	binary.Write(sigOutput, binary.BigEndian, signature)
}

// private
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

func privateKeyFromMnemonic(mnemonic string, algorithm string) crypto.PrivKey {
	seed, err := bip39.MnemonicToSeedWithErrChecking(mnemonic)
	if err != nil {
		panic(err)
	}
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, _ := hd.DerivePrivateKeyForPath(masterPriv, ch, hd.FullFundraiserPath)

	switch algorithm {
	case secpString:
		return secp256k1.PrivKeySecp256k1(derivedPriv)
	case edString:
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

	pub := pkey.PubKey()
	if !pub.VerifyBytes(signBytes, signed) {
		panic("invalid signature")
	}

	return sig
}

func getInputArgs() (args args) {
	mnemonicPtr := flag.String("mnemonic", "", "secret phrase")
	algPtr := flag.String("algorithm", "secp", "private key algorithm")
	noTrimPtr := flag.Bool("notrim", false, "don't trim sign bytes")
	flag.Parse()

	args.mnemonic = *mnemonicPtr
	args.mnemonic = strings.Trim(args.mnemonic, " ")
	validateMnemonic(args.mnemonic)

	args.algorithm = *algPtr
	args.noTrim = *noTrimPtr
	return
}
