package utils

import (
	"fmt"
)

type ConstsMap map[string]string

type Consts struct {
	STAGE     string
	ConstsMap ConstsMap
}

func NewConsts(STAGE string, ConstsMap ConstsMap) *Consts {
	return &Consts{
		ConstsMap: ConstsMap,
		STAGE:     STAGE,
	}
}

func (s *Consts) GetConst(key string) string {
	var v string
	var ok bool

	if s.STAGE == "" {
		panic("No STAGE set")
	}

	v, ok = s.ConstsMap[key]

	if !ok {
		panic(fmt.Sprintf("Variable not set: %v", key))
	}
	return v
}
