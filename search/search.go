package search

func getItemOrZero(key string, value map[string]int) int {
	if _, ok := value[key]; ok {
		return value[key]
	} else {
		return 0
	}
}

func Tf(t string, d map[string]int, tc int) float64{
	a := float64(getItemOrZero(t, d))
	b := float64(tc)

	return a / b
}