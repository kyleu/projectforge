package numeric_test

import (
	"math/rand/v2"
	"testing"

	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/app/util/numeric"
)

func randFloat64() float64 {
	//nolint:gosec // Using weak random is acceptable for benchmarks
	return rand.Float64() * 1000
}

// Benchmark functions for `numeric.Numeric` operations.
func BenchmarkNumericAdd(b *testing.B) {
	nums := make([]numeric.Numeric, 1000)
	for i := range nums {
		util.RandomFloat(1)
		nums[i] = numeric.FromFloat(randFloat64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1.Add(n2)
	}
}

func BenchmarkNumericSubtract(b *testing.B) {
	nums := make([]numeric.Numeric, 1000)
	for i := range nums {
		nums[i] = numeric.FromFloat(randFloat64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1.Subtract(n2)
	}
}

func BenchmarkNumericMultiply(b *testing.B) {
	nums := make([]numeric.Numeric, 1000)
	for i := range nums {
		nums[i] = numeric.FromFloat(randFloat64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1.Multiply(n2)
	}
}

func BenchmarkNumericDivide(b *testing.B) {
	nums := make([]numeric.Numeric, 1000)
	for i := range nums {
		nums[i] = numeric.FromFloat(randFloat64() + 1) // add 1 to avoid division by zero
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1.Divide(n2)
	}
}

// Benchmark functions for standard `float64` operations.
func BenchmarkFloat64Add(b *testing.B) {
	nums := make([]float64, 1000)
	for i := range nums {
		nums[i] = randFloat64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1 + n2
	}
}

func BenchmarkFloat64Subtract(b *testing.B) {
	nums := make([]float64, 1000)
	for i := range nums {
		nums[i] = randFloat64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1 - n2
	}
}

func BenchmarkFloat64Multiply(b *testing.B) {
	nums := make([]float64, 1000)
	for i := range nums {
		nums[i] = randFloat64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1 * n2
	}
}

func BenchmarkFloat64Divide(b *testing.B) {
	nums := make([]float64, 1000)
	for i := range nums {
		nums[i] = randFloat64() + 1 // Add 1 to avoid division by zero
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n1 := nums[i%len(nums)]
		n2 := nums[(i+1)%len(nums)]
		_ = n1 / n2
	}
}
