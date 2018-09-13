package NeuralNetwork

import "math"

// 工具类

func makeMatrix(row, column int) [][]float64{
	matrix := make([][]float64, row)
	for i := 0; i<row; i++{
		matrix[i] = make([]float64, column)
	}
	return matrix
}

func SigmodFunction(x float64) float64{
	return 1 / (1+ math.Exp(-x))
}
