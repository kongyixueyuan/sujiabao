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



type Wallet struct {
	PrivateKey ecdsa.PrivateKey			//private key struct
	PublicKey  []byte				//public key
}

func IsValidForAdress(adress []byte) bool {

	//1+20+4
	version_public_checksumBytes := Base58Decode(adress)

	checkBytes := CheckSum(version_public_checksumBytes[:len(version_public_checksumBytes) - addressCheckSumLen])

	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - addressCheckSumLen:]

	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}

	return false
}

func Ripemd160Hash(publicKey []byte) []byte {

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

func CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressCheckSumLen]
}

func (w *Wallet) GetAddress() []byte  {

	ripemd160Hash := Ripemd160Hash(w.PublicKey)

	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)

	checkSumBytes := CheckSum(version_ripemd160Hash)

	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Base58Encode(bytes)
}


func newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}




func NewWallet() *Wallet {

	privateKey,publicKey := newKeyPair()

	return &Wallet{privateKey,publicKey}
}
