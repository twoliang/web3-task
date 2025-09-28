package longestCommonPrefix

/*
编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。
*/

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		prefix = prefix[:j] //这里可能出现空字符串，比如第一个字符就不相同

		// 如果前缀已经为空，可以提前终止
		if prefix == "" {
			return ""
		}
	}

	return prefix
}
