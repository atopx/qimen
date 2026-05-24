package qimen

import (
	"errors"

	"github.com/atopx/qimen/internal/compute"
	"github.com/atopx/qimen/plate"
)

// Sentinel errors. All are errors.Is-friendly and can be wrapped freely.
var (
	// ErrUnsupportedMethod is returned when an unimplemented enum.Method
	// (MethodDay/Month/Year) is supplied via WithMethod.
	ErrUnsupportedMethod = compute.ErrUnsupportedMethod

	// ErrUnsupportedStyle is returned when StyleFly or StyleSiZhu is
	// supplied via WithStyle.
	ErrUnsupportedStyle = plate.ErrUnsupportedStyle

	// ErrInvalidTime is returned for invalid SolarTime inputs.
	ErrInvalidTime = errors.New("qimen: invalid time")
)
