package qimen

import (
	"errors"

	"github.com/atopx/qimen/almanac"
)

// Sentinel errors. All are errors.Is-friendly and can be wrapped freely.
//
// Each sentinel lives in the package that produces it: Method / Style
// are validated only at the Chart entry points, so their sentinels are
// defined here; time validation happens in the almanac layer, so
// ErrInvalidTime aliases the almanac sentinel.
var (
	// ErrUnsupportedMethod is returned when an unimplemented enum.Method
	// (MethodDay/Month/Year) is supplied via WithMethod.
	ErrUnsupportedMethod = errors.New("qimen: unsupported method")

	// ErrUnsupportedStyle is returned when StyleFly or StyleSiZhu is
	// supplied via WithStyle.
	ErrUnsupportedStyle = errors.New("qimen: unsupported style")

	// ErrInvalidTime is returned for invalid solar-time inputs. It is
	// the same sentinel as [almanac.ErrInvalidTime], so errors.Is
	// matches against either name.
	ErrInvalidTime = almanac.ErrInvalidTime
)
