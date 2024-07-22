package urn

import (
	"fmt"
	"github.com/google/uuid"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) NewURN() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return ``, fmt.Errorf("failed to generate UUID - %w", err)
	}

	if id == uuid.Nil {
		return ``, fmt.Errorf("received a nil UUID")
	}

	return `urn:uuid:` + id.String(), nil
}

func (g *Generator) Validate(urn string) bool {
	return true
}
