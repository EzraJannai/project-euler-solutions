package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

/*
Problem 84: Monopoly Odds

The goal is to find the three most frequently visited squares on a Monopoly board using two 4-sided dice. The game has special rules for "Go to Jail," "Community Chest," and "Chance" cards, which affect movement probabilities.

Solution:
This solution simulates player movement around the Monopoly board using random dice rolls and rules for special squares (Jail, Community Chest, Chance). It tracks the frequency of landing on each square over multiple iterations. The probabilities are normalized, and a PageRank algorithm is applied to determine the most visited squares.
*/

func contains(value int, list []int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func randChoose(options []int) int {
	return options[rand.Intn(len(options))]
}

func randomAlgorithm(input int) int {
	current_position := input
	double := true
	sum := 0

	first := randChoose([]int{1, 2, 3, 4})
	second := randChoose([]int{1, 2, 3, 4})

	for double {
		if first == second {
			if sum == 2 {
				current_position = 10
				return current_position
			}
			sum += 1
		} else {
			double = false
		}

		current_position = (current_position + second + first) % 40

		if contains(current_position, []int{2, 17, 33}) && rand.Intn(16) < 2 {
			current_position = randChoose([]int{0, 10})
		}

		if contains(current_position, []int{7, 22, 26}) && rand.Intn(16) < 10 {
			next_r := (current_position + 5) / 10
			next_u := 0
			if current_position == 22 {
				next_u = 28
			} else {
				next_u = 12
			}
			options := []int{0, 10, 11, 39, 24, 5, next_r, next_r, next_u, current_position - 3}
			current_position = randChoose(options)
		}

		if current_position == 30 {
			current_position = 10
			return current_position
		}

		if current_position == 10 {
			return current_position
		}

		first = randChoose([]int{1, 2, 3, 4})
		second = randChoose([]int{1, 2, 3, 4})
	}

	return current_position
}

func normalizeTable(table [][]int) [][]float64 {
	normalized := make([][]float64, len(table))
	for i, row := range table {
		normalized[i] = make([]float64, len(row))
		rowSum := 0
		for _, count := range row {
			rowSum += count
		}
		if rowSum > 0 {
			for j, count := range row {
				normalized[i][j] = float64(count) / float64(rowSum)
			}
		}
	}
	return normalized
}

func pageRank(matrix [][]float64, iterations int, damping float64) []float64 {
	n := len(matrix)
	rank := make([]float64, n)
	for i := range rank {
		rank[i] = 1.0 / float64(n)
	}

	for iter := 0; iter < iterations; iter++ {
		newRank := make([]float64, n)
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < n; j++ {
				sum += matrix[j][i] * rank[j]
			}
			newRank[i] = (1-damping)/float64(n) + damping*sum
		}
		rank = newRank
	}
	return rank
}

func main() {
	rand.Seed(time.Now().UnixNano())

	resultTable := make([][]int, 40)
	for i := range resultTable {
		resultTable[i] = make([]int, 40)
	}

	for startPos := 0; startPos < 40; startPos++ {
		for sample := 0; sample < 1000; sample++ {
			finalPosition := randomAlgorithm(startPos)
			resultTable[startPos][finalPosition]++
		}
	}

	normalizedTable := normalizeTable(resultTable)
	pageRankValues := pageRank(normalizedTable, 1000, 0.85)

	type squareRank struct {
		square int
		rank   float64
	}

	ranks := make([]squareRank, len(pageRankValues))
	for i, rankValue := range pageRankValues {
		ranks[i] = squareRank{i, rankValue}
	}

	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].rank > ranks[j].rank
	})

	fmt.Println("Square Rankings (most visited to least):")
	for i, square := range ranks {
		fmt.Printf("%2d: Square %2d with PageRank %.4f\n", i+1, square.square, square.rank)
	}
}
