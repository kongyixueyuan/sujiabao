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

const walletFile  = "Wallets_%s.dat"

type SJB_Wallets struct {
	SJB_WalletsMap map[string]*SJB_Wallet
}


func SJB_NewWallets(nodeId string) (*SJB_Wallets,error){

	walletFileName := fmt.Sprintf(walletFile,nodeId)
	if _, err := os.Stat(walletFileName); os.IsNotExist(err) {
		wallets := &SJB_Wallets{}
		wallets.SJB_WalletsMap = make(map[string]*SJB_Wallet)
		return wallets,err
	}

	//else load and return it
	fileContent, err := ioutil.ReadFile(walletFileName)
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



func (w *SJB_Wallets) SJB_SaveWallets(nodeId string)  {
	var content bytes.Buffer
	walletFileName := fmt.Sprintf(walletFile,nodeId)
	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	//write
	err = ioutil.WriteFile(walletFileName, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}


}


func (w *SJB_Wallets) SJB_CreateNewWallet(nodeId string)  {

	wallet := SJB_NewWallet()
	fmt.Printf("new address：%s\n",wallet.SJB_GetAddress())
	w.SJB_WalletsMap[string(wallet.SJB_GetAddress())] = wallet
	w.SJB_SaveWallets(nodeId)
}

