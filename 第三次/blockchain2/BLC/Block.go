package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
	"fmt"
)

type Block struct {
	//1. 区块高度
	Height int64
	//2. 上一个区块HASH
	PrevBlockHash []byte
	//3. 交易数据
	Data []byte
	//4. 时间戳
	Timestamp int64
	//5. Hash
	Hash []byte
	// 6. Nonce
	Nonce int64
}


func (block *Block)Serialize() []byte{
	var result  bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err);
	}

	return result.Bytes();
}

func DeSerianlize(blockBytes []byte) *Block{

	var block Block;
	decoder := gob.NewDecoder(bytes.NewBuffer(blockBytes))
	err := decoder.Decode(&blockBytes)
	if err != nil{
		log.Panic(err)
	}
	return &block
}


func NewBlock(data string, height int64,preBlockHash []byte) *Block{


	//new block
	block := &Block{height,preBlockHash,[]byte(data),time.Now().Unix(),nil,0}

	pow := NewProofOfWork(block);

	hash,nonce := pow.Run();

	block.Hash = hash[:]
	block.Nonce = nonce;

	fmt.Printf("*******第%d区块生成成功*********",height);
	fmt.Println()
	return block;
}



func CreateGenesisBlock(data string) *Block{

	return NewBlock(data,1, []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}

