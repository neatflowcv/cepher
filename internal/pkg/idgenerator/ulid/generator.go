package ulid

import "github.com/oklog/ulid/v2"

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateID() string {
	return ulid.Make().String()
}
