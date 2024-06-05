package util

import "encoding/json"

func TwoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for idx, num := range nums {

		if v, found := m[target-num]; found {
			return []int{v, idx}
		}

		m[num] = idx
	}
	return nil
}

func JSONString(data interface{}) (string, error) {
	databyte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(databyte), nil
}
