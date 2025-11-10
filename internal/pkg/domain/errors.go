package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

func InvalidParameterError(param string) error {
	return fmt.Errorf("invalid %s: %w", param, ErrInvalidParameter)
}
