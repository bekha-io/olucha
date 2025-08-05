package api

type DSLParser interface {
	Parse(raw []byte) (DSL, error)
}