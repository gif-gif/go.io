package bean

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
