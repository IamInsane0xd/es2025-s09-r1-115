package models

type Container struct {
	Id        string "json:\"id\""
	BlockId   int64  "json:\"blockId\""
	BayNum    int64  "json:\"bayNum\""
	StackNum  int64  "json:\"stackNum\""
	TierNum   int64  "json:\"tierNum\""
	ArrivedAt string "json:\"arrivedAt\""
}

func NewContainer(id string, blockId int64, bayNum int64, stackNum int64, tierNum int64, arrivedAt string) *Container {
	return &Container{
		id,
		blockId,
		bayNum,
		stackNum,
		tierNum,
		arrivedAt,
	}
}
