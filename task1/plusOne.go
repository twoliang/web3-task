package task1

/*
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。
这个大整数不包含任何前导 0。
将大整数加 1，并返回结果的数字数组。
*/

func plusOne(digits []int) []int {
	if len(digits) == 1 && digits[0] == 9 {
		return []int{1, 0}
	}

	tail := digits[len(digits)-1]
	if tail < 9 {
		digits[len(digits)-1] = tail + 1
		return digits
	}
	digits[len(digits)-2] = digits[len(digits)-2] + 1
	digits[len(digits)-1] = 0
	return digits
}

// 完善版本
func plusOne1(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	// 如果所有位都是9，需要在前面加1
	return append([]int{1}, digits...)
}
