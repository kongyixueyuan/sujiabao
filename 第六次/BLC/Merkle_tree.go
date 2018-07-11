package BLC

import (
	"crypto/sha256"
)

type SJB_MerkleTree struct{
	SJB_RootNode *SJB_MerkleNode
}

type SJB_MerkleNode struct{
	SJB_Leftnode *SJB_MerkleNode
	SJB_Rightnode *SJB_MerkleNode
	SJB_Data []byte
}


func SJB_NewMerkleTree(data [][]byte) *SJB_MerkleTree{

	var nodes  []SJB_MerkleNode

	if(len(data)%2 != 0){
		data = append(data,data[len(data)-1])
	}

	for _,nodedata := range data{
		newnode := SJB_NewMerkleNode(nil,nil,nodedata)
		nodes = append(nodes,*newnode)
	}

	var newNodeLever []SJB_MerkleNode
	for i := 0; i < len(data)/2; i++  {
		for j := 0; j<len(nodes);j+=2{
			node := SJB_NewMerkleNode(&nodes[j], &nodes[j+1],nil)
			newNodeLever = append(newNodeLever,*node)
		}
		nodes = newNodeLever
	}

	newMerkleTree := SJB_MerkleTree{&nodes[0]}

	return &newMerkleTree
}

func SJB_NewMerkleNode(left,right *SJB_MerkleNode, data []byte) *SJB_MerkleNode{

	newnode := SJB_MerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		newnode.SJB_Data = hash[:]
	}else{
		perhash := append(left.SJB_Data,right.SJB_Data...)
		hash := sha256.Sum256(perhash)
		newnode.SJB_Data = hash[:]
	}
	newnode.SJB_Leftnode = left
	newnode.SJB_Rightnode = right

	return &newnode
}

