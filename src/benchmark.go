package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func BenchmarkContains(b *testing.B) {
	url := "https://www.xxx.com/sg-en/laptops/16/RZ09/awdiub/awdda"
	for i := 0; i < b.N; i++ {
		if strings.Contains(url, "RZ") || strings.Contains(url, "RS") || strings.Contains(url, "RC") {
			// Do something
		}
	}
}

func BenchmarkRegex(b *testing.B) {
	url := "https://www.xxx.com/sg-en/laptops/16/RZ09/awdiub/awdda"
	re := regexp.MustCompile(`R[ZSC]`)
	for i := 0; i < b.N; i++ {
		if re.MatchString(url) {
			// Do something
		}
	}
}

func testBenchMark() {
	fmt.Println("Benchmarking extractProductCode with Contains...")
	containsResult := testing.Benchmark(BenchmarkContains)
	fmt.Println(containsResult)

	fmt.Println("Benchmarking extractProductCode with Regex...")
	regexResult := testing.Benchmark(BenchmarkRegex)
	fmt.Println(regexResult)
}
