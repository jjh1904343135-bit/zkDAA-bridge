package main

import (
	"fmt"
	"log"

	"crypto/ecdsa"

	crypto "github.com/ethereum/go-ethereum/crypto"
)

func main() {
	bgPrivKeyHex := "a1dc2753c6ec7cf4cd507eecd81072db11f9f26189f2ab10923c963343b92e24" // 换成你自己的

	priv, err := crypto.HexToECDSA(bgPrivKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	pub := priv.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pub)

	fmt.Println("背景压测账户地址:", addr.Hex())
}
