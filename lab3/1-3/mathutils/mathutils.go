package mathutils

func Factorial(n int) uint64 {

	var result uint64 = 1
	for i := 1; i <= n; i++ {
		result *= uint64(i)
	}
	return result
}
