package godb

type NumberRangeController[T int64 | int | float64] struct {
	Min *T `json:"min,optional"`
	Max *T `json:"max,optional"`
}

func (c *NumberRangeController[T]) Where(column string) (string, []any) {
	if c.Min == nil && c.Max == nil {
		return "", []any{}
	}

	var whereString string
	args := make([]any, 0, 2)
	if c.Min != nil {
		whereString = " " + column + "  >= ? "
		args = append(args, *c.Min)
	}
	if c.Max != nil {
		if len(whereString) > 0 {
			whereString += " and " + column + "  < ? "
		} else {
			whereString = " " + column + "  < ? "
		}
		args = append(args, *c.Max)
	}

	return whereString, args
}
