package BLC

type TXOutput struct {
	Value int64
	ScriptPubKey string
}

func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {

	return txOutput.ScriptPubKey == address
}


