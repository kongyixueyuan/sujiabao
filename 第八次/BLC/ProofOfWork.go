package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)



const targetBit  = 16



type SJB_ProofOfWork struct {
	SJB_Block *SJB_Block
	SJB_target *big.Int
}


func (pow *SJB_ProofOfWork) SJB_SeriesData(nonce int) []byte{
	data := bytes.Join(
		[][]byte{
			pow.SJB_Block.SJB_PrevBlockHash,
			pow.SJB_Block.SJB_HashTransactions(),
			SJB_IntToHex(pow.SJB_Block.SJB_Timestamp),
			SJB_IntToHex(int64(targetBit)),
			SJB_IntToHex(int64(nonce)),
			SJB_IntToHex(int64(pow.SJB_Block.SJB_Height)),
		},
		[]byte{},
	)

	return data
}


func (pofwork *SJB_ProofOfWork)SJB_Run() ([]byte,int64){
	nonce := 0;
	var hashInt big.Int
	var hash [32]byte

	for {
		dataBytes := pofwork.SJB_SeriesData(nonce)
		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x",hash)
		hashInt.SetBytes(hash[:])
		if pofwork.SJB_target.Cmp(&hashInt) == 1{
			break
		}
		nonce = nonce +1

	}

	return hash[:],int64(nonce)
}

func SJB_NewProofOfWork(block *SJB_Block) *SJB_ProofOfWork  {

	target := big.NewInt(1)
	target = target.Lsh(target,256-targetBit)
	return &SJB_ProofOfWork{block,target}
}






