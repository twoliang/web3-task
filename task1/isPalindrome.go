package isPalindrome

import "strconv"

/*
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。

回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。

例如，121 是回文，而 123 不是。
*/

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	s := strconv.Itoa(x)
	list := []rune(s)

	head, tail := 0, len(list)-1

	for head < tail {
		if list[head] != list[tail] {
			return false
		}
		head++
		tail--
	}

	return true
}
