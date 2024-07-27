package main

import "fmt"

type Printer interface {
	printAlignment(alignment Alignment)

	printMatrix(matrix [][]StepState)

	printScore(matrix [][]StepState)

	printSequences(seqA, seqB string)
}

type SimplePrinter struct {
}

func (p SimplePrinter) printMatrix(matrix [][]StepState) {

	if len(matrix) > 20 {
		fmt.Println("Matrix too large, skipping printing matrix")
		return
	}
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

func (p SimplePrinter) printSequences(seqA, seqB string) {
	fmt.Printf("\nSequence 1: %s\nSequence 2: %s\n", seqA, seqB)
}

func NewPrinter() Printer {
	return SimplePrinter{}
}

func NoopPrinter() Printer {
	return NullPrinter{}
}

type NullPrinter struct {
}

func (p NullPrinter) printAlignment(alignment Alignment) {

}

func (p NullPrinter) printMatrix(matrix [][]StepState) {

}

func (p NullPrinter) printScore(matrix [][]StepState) {

}

func (p NullPrinter) printSequences(seqA, seqB string) {

}
