package pk

import (
	"fmt"
)
// ID represents a type alias for a primary key
type ID int

// Validate returns error if the ID value is not right 
func (id ID) Validate() error {
	if id < 1 {
		return fmt.Errorf("pk.ID is not a positive number")
	}
	return nil
}