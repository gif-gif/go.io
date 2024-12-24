package bean

const (
	OpEquals         = "equals"
	OpLessOrEquals   = "lessEqual"
	OpGreaterOrEqual = "greaterEqual"
	OpGreater        = "greater"
	OpLess           = "less"
)

type NumberRange[T any] struct {
	Min *T `json:"min,optional"`
	Max *T `json:"max,optional"`
}

type Int64Range struct {
	Min *int64 `json:"min,optional"`
	Max *int64 `json:"max,optional"`
}

type Float64Range struct {
	Min *float64 `json:"min,optional"`
	Max *float64 `json:"max,optional"`
}

type Int64Method struct {
	Method string `json:"method,optional"`
	Value  int64  `json:"value,optional"`
}

func (r *Int64Method) Check(val int64) bool {
	return CheckValue(r, val)
}
func (r *Int64Method) GetValue() int64 {
	return r.Value
}

func (r *Int64Method) GetMethod() string {
	return r.Method
}

type IOperationMethod[T any] interface {
	GetMethod() string
	GetValue() T
}

func CheckValue(r IOperationMethod[int64], val int64) bool {
	method := r.GetMethod()
	value := r.GetValue()
	switch method {
	case OpEquals:
		return val == value
	case OpLessOrEquals:
		return val <= value
	case OpGreaterOrEqual:
		return val >= value
	case OpGreater:
		return val > value
	case OpLess:
		return val < value
	}
	return false
}
