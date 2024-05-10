package goutils

func SplitStringArray(arr []string, size int) (list [][]string) {
	l := len(arr)

	if l == 0 {
		list = make([][]string, 0)
		return
	}

	if l < size {
		list = [][]string{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitIntArray(arr []int, size int) (list [][]int) {
	l := len(arr)

	if l == 0 {
		list = make([][]int, 0)
		return
	}

	if l < size {
		list = [][]int{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitInt64Array(arr []int64, size int) (list [][]int64) {
	l := len(arr)

	if l == 0 {
		list = make([][]int64, 0)
		return
	}

	if l < size {
		list = [][]int64{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitArray(arr []interface{}, size int) (list [][]interface{}) {
	l := len(arr)

	if l == 0 {
		list = make([][]interface{}, 0)
		return
	}

	if l < size {
		list = [][]interface{}{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}
