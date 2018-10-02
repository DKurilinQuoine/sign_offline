# SignOffline
This utility allows signing transactions for offline users made by basecli.

## Build
make - to build on the local machine
make test - will produce a jenkins compatible report in the build folder
make docker - builds a docker container, runs tests in it and stores the report in the build folder(made for Jenkins)

## How to use

./sign_offline -mnemonic="string" [-algorithm=("secp"|"ed")] [-notrim] < bytes_to_sign 2> signature_output

1. mnemonic - mandatory flag. Contains a mnemonic string. Must be 24 words long separated with a space symbol.
2. algorithm ("secp"|"ed") - optional flag. Specifies which algorithm to use to generate a private key from a mnemonic. Default: "secp"
3. notrim - optional flag. Turns off trimming of new line symbols and spaces from the bytes_to_sign. Default: false
4. bytes_to_sign - A file with bytes that has to be signed.
5. signature_output - An output file with a signature. 
#### ! The utility outputs a signature to the stderr.
