package models

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

var (
	NotFoundErr      = errors.New("error: container not found")
	AlreadyExistsErr = errors.New("error: container already exists")
	IncorrectPosErr  = errors.New("error: incorrect container position")
	ParamError       = errors.New("error: incorrect parameter value")
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

func (c *ContStore) Get(id string) (Container, error) {
	for _, val := range c.list {
		if val.Id == id {
			return val, nil
		}
	}

	return Container{}, NotFoundErr
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

func (c *ContStore) GetByQuery(query url.Values) ([]Container, error) {
	var containers []Container

	if query.Has("id") {
		if len(query["id"][0]) == 11 {
			container, err := c.Get(query["id"][0])
			return []Container{container}, err
		}

		for _, val := range c.list {
			if !strings.Contains(val.Id, query["id"][0]) {
				continue
			}

			containers = append(containers, val)
		}

		if len(containers) == 0 {
			return []Container{}, NotFoundErr
		}

		return containers, nil
	}

	var blockId int64
	var bayNum int64
	var stackNum int64
	var tierNum int64
	var err error

	if query.Has("blockId") {
		if blockId, err = strconv.ParseInt(query["blockId"][0], 10, 64); err != nil {
			return []Container{}, ParamError
		}
	}

	if query.Has("bayNum") {
		if bayNum, err = strconv.ParseInt(query["bayNum"][0], 10, 64); err != nil {
			return []Container{}, ParamError
		}
	}

	if query.Has("stackNum") {
		if stackNum, err = strconv.ParseInt(query["stackNum"][0], 10, 64); err != nil {
			return []Container{}, ParamError
		}
	}

	if query.Has("tierNum") {
		if tierNum, err = strconv.ParseInt(query["tierNum"][0], 10, 64); err != nil {
			return []Container{}, ParamError
		}
	}

	for _, val := range c.list {
		if query.Has("blockId") && (val.BlockId != blockId) || query.Has("bayNum") && (val.BayNum != bayNum) ||
			query.Has("stackNum") && (val.StackNum != stackNum) || query.Has("tierNum") && (val.TierNum != tierNum) {
			continue
		}

		containers = append(containers, val)
	}

	if len(containers) == 0 {
		return []Container{}, NotFoundErr
	}

	return containers, nil
}
