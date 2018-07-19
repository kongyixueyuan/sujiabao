package BLC


var knowNodes = []string{"localhost:3000"}
var nodeAddress string
var transactionArray [][]byte
var MyminerAddress string
var mempool = make(map[string]SJB_Transaction)

func SJB_nodeIsKnow(address string) bool{
	for _,add := range knowNodes{
		if add == address{
			return true
		}
	}
	return false
}