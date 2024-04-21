package main

import "fmt"

type Printer interface {
	printAlignment(alignment Alignment)

	printMatrix(matrix [][]StepState)

	printScore(matrix [][]StepState)
}

type SimplePrinter struct {
}

func (p SimplePrinter) printMatrix(matrix [][]StepState) {
	fmt.Print("------------------------------\n")
	for _, row := range matrix {
		for _, col := range row {
			fmt.Printf("|%3d ", col.Score)
		}
		fmt.Print("|\n")
	}
	fmt.Print("------------------------------\n")

}

func (p SimplePrinter) printAlignment(alignment Alignment) {
	fmt.Printf("\n%s\n%s\n%s\n", alignment.Sequence1, alignment.Comparison, alignment.Sequence2)
}

func (p SimplePrinter) printScore(matrix [][]StepState) {
	fmt.Printf("Max score : %d\n", matrix[len(matrix)-1][len(matrix[0])-1])
}

func NewPrinter() Printer {
	return SimplePrinter{}
}
