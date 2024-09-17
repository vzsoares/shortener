package utils

import (
	"fmt"
)

type ConstsMap map[string]string

type Consts struct {
	STAGE   string
	ProdMap ConstsMap
	DevMap  ConstsMap
}

func NewConsts(STAGE string, ProdMap ConstsMap, DevMap ConstsMap) *Consts {
	return &Consts{
		ProdMap: ProdMap,
		DevMap:  DevMap,
		STAGE:   STAGE,
	}
}

func (s *Consts) GetConst(key string) string {
	var v string
	var ok bool

	if s.STAGE == "" {
		panic("No STAGE set")
	}

	if s.STAGE == "dev" {
		v, ok = s.DevMap[key]
	} else {
		v, ok = s.ProdMap[key]
	}

	if !ok {
		panic(fmt.Sprintf("Variable not set: %v", key))
	}
	return v
}
