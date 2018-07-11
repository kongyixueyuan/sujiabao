package BLC


func (cli *SJB_CLI) SJB_createGenesisBlockchain(address string)  {

	blockchain := SJB_CreateBlockchainWithGenesisBlock(address)

	defer blockchain.SJB_DB.Close()
	utxoSet := &SJB_UTXOSet{blockchain}
	utxoSet.SJB_ResetUTXOSet()
}