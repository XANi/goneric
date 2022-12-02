package goneric

func Sum[T Number](n ...T) (sum T) {
	for _, v := range n {
		sum = sum + v
	}
	return sum
}
