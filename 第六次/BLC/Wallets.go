package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"fmt"
)

const walletFile  = "AddressWallets.dat"

type SJB_Wallets struct {
	SJB_WalletsMap map[string]*SJB_Wallet
}



func SJB_NewWallets() (*SJB_Wallets,error){

	//if wallets data file not exist ,creat one
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &SJB_Wallets{}
		wallets.SJB_WalletsMap = make(map[string]*SJB_Wallet)
		return wallets,err
	}

	//else load and return it
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}

	var wallets SJB_Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}

	return &wallets,nil
}



func (w *SJB_Wallets) SJB_SaveWallets()  {
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


func (w *SJB_Wallets) SJB_CreateNewWallet()  {

	wallet := SJB_NewWallet()
	fmt.Printf("new address：%s\n",wallet.SJB_GetAddress())
	w.SJB_WalletsMap[string(wallet.SJB_GetAddress())] = wallet
	w.SJB_SaveWallets()
}

