package uuid_generator

import "github.com/google/uuid"

type Generator struct {
}

func (g *Generator) Generate() string {
	return uuid.NewString()
}

func New() *Generator {
	return &Generator{}
}
