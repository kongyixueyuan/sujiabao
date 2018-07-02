package BLC

type TXInput struct {
	TxHash      []byte
	Vout      int
	ScriptSig string
}

func (txInput *TXInput) UnLockWithAddress(address string) bool {

	return txInput.ScriptSig == address
}