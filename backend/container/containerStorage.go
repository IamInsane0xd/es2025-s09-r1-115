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

func (c *ContStore) GetByBlockId(blockId int) ([]Container, error) {
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

func (c *ContStore) List() map[string]Container {
	return c.list
}
