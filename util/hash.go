package util

type F func(string) uint64 //创建一个函数类型
//把所有hash函数放入切片中
func NewFunc() []F {
	m := make([]F, 0)
	var f F
	f = BKDRHash
	m = append(m, f)
	f = SDBMHash
	m = append(m, f)
	f = DJBHash
	m = append(m, f)
	return m
}
func BKDRHash(str string) uint64 {
	seed := uint64(131) // 31 131 1313 13131 131313 etc..
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + uint64(str[i])
	}
	return hash & 0x7FFFFFFF
}
func SDBMHash(str string) uint64 {
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = uint64(str[i]) + (hash << 6) + (hash << 16) - hash
	}
	return hash & 0x7FFFFFFF
}
func DJBHash(str string) uint64 {
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = ((hash << 5) + hash) + uint64(str[i])
	}
	return hash & 0x7FFFFFFF
}
