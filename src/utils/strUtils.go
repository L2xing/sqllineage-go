package utils

import "strconv"

// int64转字符串
func Int642str(i int64) string {
	return strconv.FormatInt(i, 10)
}

// int转字符串
func Int2str(i int) string {
	return Int642str(int64(i))
}

// uint64转字符串
func Uint642str(i uint64) string {
	return strconv.FormatUint(i, 10)
}

// uint64转字符串
func Uint2str(i uint) string {
	return Uint642str(uint64(i))
}
