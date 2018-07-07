package BLC

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"crypto/elliptic"
	"io/ioutil"
	"log"
	"os"
)

const walletFile  = "BlockChainWallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}


func NewWallets() (*Wallets,error){

	//if wallets data file not exist ,creat one
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletsMap = make(map[string]*Wallet)
		return wallets,err
	}

	//else load and return it
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}


func (w *Wallets) CreateNewWallet()  {

	wallet := NewWallet()
	fmt.Printf("new address：%s\n",wallet.GetAddress())
	w.WalletsMap[string(wallet.GetAddress())] = wallet
	w.SaveWallets()
}

func (w *Wallets) SaveWallets()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	//write
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}


}

