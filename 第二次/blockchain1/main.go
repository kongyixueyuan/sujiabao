package main

import (
	"PublicChain1801/苏家宝/第二次/blockchain1/BLC"
)

func main() {
	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	// 新区块
	blockchain.AddBlockToBlockchain("我是第二个", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	blockchain.AddBlockToBlockchain("我是第三个", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	blockchain.AddBlockToBlockchain("我是第四个", blockchain.Blocks[len(blockchain.Blocks)-1].Height+1, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

}
