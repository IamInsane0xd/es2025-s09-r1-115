package container

import (
	"errors"
)

var (
	NotFoundErr      = errors.New("error: container not found")
	AlreadyExistsErr = errors.New("error: container already exists")
	IncorrectPosErr  = errors.New("error: incorrect container position")
)

type ContStore struct {
	list map[string]Container
}

func NewContainerStore() *ContStore {
	list := make(map[string]Container)

	return &ContStore{
		list,
	}
}

func (c *ContStore) Add(name string, container Container) error {
	for _, cont := range c.list {
		if (cont.BlockId == container.BlockId && cont.BayNum == container.BayNum &&
			cont.StackNum == container.StackNum && cont.TierNum == container.TierNum) ||
			(1 > container.BlockId || container.BlockId > 4 ||
				1 > container.BayNum || container.BayNum > 5 ||
				1 > container.StackNum || container.StackNum > 5 ||
				1 > container.TierNum || container.TierNum > 5) {
			return IncorrectPosErr
		}

		if cont.Id == container.Id {
			return AlreadyExistsErr
		}
	}

	c.list[name] = container
	return nil
}

func (c *ContStore) GetBlock(blockId int64) ([]Container, error) {
	var containers []Container

	for _, val := range c.list {
		if val.BlockId == blockId {
			containers = append(containers, val)
		}
	}

	if len(containers) == 0 {
		return []Container{}, NotFoundErr
	}

	return containers, nil
}

func (c *ContStore) GetBay(blockId int64, bayNum int64) ([]Container, error) {
	block, err := c.GetBlock(blockId)

	if err != nil {
		return []Container{}, err
	}

	var containers []Container

	for _, val := range block {
		if val.BayNum == bayNum {
			containers = append(containers, val)
		}
	}

	if len(containers) == 0 {
		return []Container{}, err
	}

	return containers, nil
}

func (c *ContStore) GetStack(blockId int64, bayNum int64, stackNum int64) ([]Container, error) {
	bay, err := c.GetBay(blockId, bayNum)

	if err != nil {
		return []Container{}, err
	}

	var containers []Container

	for _, val := range bay {
		if val.StackNum == stackNum {
			containers = append(containers, val)
		}
	}

	if len(containers) == 0 {
		return []Container{}, NotFoundErr
	}

	return containers, nil
}

func (c *ContStore) GetContainer(blockId int64, bayNum int64, stackNum int64, tierNum int64) (Container, error) {
	stack, err := c.GetStack(blockId, bayNum, stackNum)

	if err != nil {
		return Container{}, err
	}

	for _, val := range stack {
		if val.TierNum == tierNum {
			return val, nil
		}
	}

	return Container{}, NotFoundErr
}

func (c *ContStore) List() map[string]Container {
	return c.list
}
