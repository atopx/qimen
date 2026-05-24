package qimen

import "fmt"

// QimenError 是奇门遁甲计算过程中可能出现的错误。
type QimenError struct {
	Code QimenErrorCode
	// 详细信息: 对 UnsupportedMethod/UnsupportedChartType 为枚举名;
	// 对 UnsupportedTerm 为节气名。
	Detail string
}

// QimenErrorCode 错误类别。
type QimenErrorCode int

const (
	// ErrCodeUnsupportedMethod 不支持的起局方法。
	ErrCodeUnsupportedMethod QimenErrorCode = iota
	// ErrCodeUnsupportedChartType 不支持的盘式。
	ErrCodeUnsupportedChartType
	// ErrCodeUnsupportedTerm 节气索引越界。
	ErrCodeUnsupportedTerm
)

func (e *QimenError) Error() string {
	switch e.Code {
	case ErrCodeUnsupportedMethod:
		return fmt.Sprintf("unsupported qimen method: %s", e.Detail)
	case ErrCodeUnsupportedChartType:
		return fmt.Sprintf("unsupported qimen chart type: %s", e.Detail)
	case ErrCodeUnsupportedTerm:
		return fmt.Sprintf("unsupported solar term for qimen: %s", e.Detail)
	default:
		return fmt.Sprintf("qimen error: %s", e.Detail)
	}
}

// newUnsupportedMethod 构造不支持方法错误。
func newUnsupportedMethod(m QimenMethod) error {
	return &QimenError{Code: ErrCodeUnsupportedMethod, Detail: m.Name()}
}

// newUnsupportedChartType 构造不支持盘式错误。
func newUnsupportedChartType(c QimenChartType) error {
	return &QimenError{Code: ErrCodeUnsupportedChartType, Detail: c.Name()}
}

// newUnsupportedTerm 构造不支持节气错误。
func newUnsupportedTerm(name string) error {
	return &QimenError{Code: ErrCodeUnsupportedTerm, Detail: name}
}
