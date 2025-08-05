package api

import (
	"encoding/json"
	"errors"
)

var (
	ErrInvalidFormat = errors.New("invalid DSL format")
)

type DSLParser interface {
	Parse(raw []byte) (DSL, error)
}

type jsonParser struct {
}

func (p *jsonParser) Parse(raw []byte) (DSL, error) {
	var dsl DSL
	err := json.Unmarshal(raw, &dsl)
	return dsl, errors.Join(err, ErrInvalidFormat)
}