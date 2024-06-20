package godb

type SelectController[T int64 | string] struct {
	Values  []T  `json:"values,optional"`
	Exclude bool `json:"exclude,optional"`
}

func (c *SelectController[T]) ClickHouseWhere(column string) (string, []T) {
	if len(c.Values) == 0 {
		return "", nil
	}

	var whereString string
	if c.Exclude {
		whereString = " not in ? "
	} else {
		whereString = " in ? "
	}

	return " " + column + " " + whereString, c.Values
}

func (c *SelectController[T]) MysqlWhere(column string) (string, []any) {
	if len(c.Values) == 0 {
		return "", nil
	}

	var whereString string
	if c.Exclude {
		whereString = column + " not in  "
	} else {
		whereString = column + " in  "
	}

	conditions, params := GenerateSliceIn[T](c.Values)
	return whereString + conditions, params
}
