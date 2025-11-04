package main

import "fmt"

type freqMap map[int]int

func countFrequency(nums []int) freqMap {
	result := make(freqMap)
	for _, n := range nums {
		result[n]++
	}
	return result
}

func main() {
	arr := []int{1, 2, 3, 2, 3, 3, 4}
	fmt.Println(countFrequency(arr))
}
