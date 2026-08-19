package models

import (
	"fmt"
)

type ActorDto struct {
	Name      string  `json:"name"`
	BirthDate string  `json:"birthDate"`
	Movies    []Movie `json:"movies"`
}

func (m ActorDto) Validate() (bool, error) {
	fmt.Errorf("Some error")
	return true, nil
}
