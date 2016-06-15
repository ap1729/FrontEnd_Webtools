package perf

func Filter(values []float64, filtIndices []uint) []float64 {
	val := make([]float64, len(filtIndices))
	for i := 0; i < len(filtIndices); i++ {
		val[i] = values[filtIndices[i]]
	}
	return val
}

func Sort(values []float64) ([]float64, []uint) {
	ind := make([]uint, len(values))
	val := make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		ind[i] = uint(i)
		val[i] = values[i]
	}
	for i := 0; i < len(val); i++ {
		for j := 1; j < len(val); j++ {
			if val[j-1] < val[j] {
				temp1 := val[j-1]
				val[j-1] = val[j]
				val[j] = temp1
				temp2 := ind[j-1]
				ind[j-1] = ind[j]
				ind[j] = temp2
			}
		}
	}
	return val, ind
}
