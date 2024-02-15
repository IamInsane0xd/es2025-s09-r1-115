package container

type Container struct {
	Id        string "json:\"id\""
	BlockId   int    "json:\"blockId\""
	BayNum    int    "json:\"bayNum\""
	StackNum  int    "json:\"stackNum\""
	TierNum   int    "json:\"tierNum\""
	ArrivedAt string "json:\"arrivedAt\""
}

func NewContainer(id string, blockId int, bayNum int, stackNum int, tierNum int, arrivedAt string) *Container {
	return &Container{
		id,
		blockId,
		bayNum,
		stackNum,
		tierNum,
		arrivedAt,
	}
}
