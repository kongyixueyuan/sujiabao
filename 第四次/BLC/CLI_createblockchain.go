package BLC


func (cli *CLI) createGenesisBlockchain(address string)  {

	blockchain := CreateBlockchainWithGenesisBlock(address)

	defer blockchain.DB.Close()
}