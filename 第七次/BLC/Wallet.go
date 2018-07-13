package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

const version = byte(0x00)
const addressCheckSumLen = 4



type SJB_Wallet struct {
	SJB_PrivateKey ecdsa.PrivateKey			//private key struct
	SJB_PublicKey  []byte				//public key
}

func SJB_IsValidForAdress(adress []byte) bool {
	//1+20+4
	if len(adress) == 0{
		return false
	}
	version_public_checksumBytes := SJB_Base58Decode(adress)
	checkBytes := SJB_CheckSum(version_public_checksumBytes[:len(version_public_checksumBytes) - addressCheckSumLen])
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - addressCheckSumLen:]
	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}
	return false
}

func SJB_Ripemd160Hash(publicKey []byte) []byte {

	//hash256

	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash_256 := hash256.Sum(nil)
	//hash160
	ripemd_160 := ripemd160.New()
	ripemd_160.Write(hash_256)
	hash_160 := ripemd_160.Sum(nil)

	return hash_160;
}

func SJB_CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressCheckSumLen]
}

func (w *SJB_Wallet) SJB_GetAddress() []byte  {

	ripemd160Hash := SJB_Ripemd160Hash(w.SJB_PublicKey)

	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)

	checkSumBytes := SJB_CheckSum(version_ripemd160Hash)

	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return SJB_Base58Encode(bytes)
}


func SJB_newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

func SJB_NewWallet() *SJB_Wallet {

	privateKey,publicKey := SJB_newKeyPair()

	return &SJB_Wallet{privateKey,publicKey}
}
