package pkg

import (
	"fmt"
	"github.com/google/uuid"
)

type URN struct{}

func NewURN() *URN {
	return &URN{}
}

func (u *URN) New() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return ``, fmt.Errorf("failed to generate UUID - %w", err)
	}

	if id == uuid.Nil {
		return ``, fmt.Errorf("received a nil UUID")
	}

	return `urn:uuid:` + id.String(), nil
}

func (u *URN) Validate(urn string) bool {
	return true
}
