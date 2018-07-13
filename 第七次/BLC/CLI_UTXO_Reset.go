package BLC


func (cli *SJB_CLI) SJB_resetUTXOSet(nodeID string)  {

	blockchain := SJB_BlockchainObject(nodeID)

	defer blockchain.SJB_DB.Close()

	utxoSet := &SJB_UTXOSet{blockchain}

	utxoSet.SJB_ResetUTXOSet()

}
