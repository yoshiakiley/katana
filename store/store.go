package store

import (
	"fmt"
)

type QueryOpType uint

type CompareType struct {
	Op    QueryOpType
	Value any
}

type ValueType struct {
	Value string
	Regex bool
}

type PaginationType struct {
	Limit int
	Skip  int
}

const (
	LT QueryOpType = iota
	LTE
	GT
	GTE
	ALL
	OR
	Ne
)

var DataNotFound = fmt.Errorf("dataNotFound")
var Store string
