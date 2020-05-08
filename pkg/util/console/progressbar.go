package console

import (
	"fmt"
	"strings"
)

type ProgressBar interface {
	Increment(inc int64)
	Done()
}

type Printer func(int64) string

func NewProgressBar(max int64, printer Printer) ProgressBar {
	cpb := consoleProgressBar{
		maxValue:       max,
		currentValue:   0,
		valuePrinter:   printer,
		maxValueString: printer(max),
	}

	fmt.Print(cpb.buildConsoleLine())
	return &cpb
}

type consoleProgressBar struct {
	maxValue     int64
	currentValue int64
	valuePrinter Printer

	maxValueString string
}

const progressBarLength = 50

func (cpb *consoleProgressBar) Increment(inc int64) {
	if cpb.currentValue+inc > cpb.maxValue {
		cpb.currentValue = cpb.maxValue
	} else {
		cpb.currentValue += inc
	}

	line := cpb.buildConsoleLine()
	fmt.Printf("\r\033[K%s", line)
}

func (cpb *consoleProgressBar) Done() {
	if cpb.currentValue < cpb.maxValue {
		cpb.Increment(cpb.maxValue - cpb.currentValue)
	}

	fmt.Println()
}

func (cpb *consoleProgressBar) buildConsoleLine() string {
	equalSigns := int(progressBarLength * cpb.currentValue / cpb.maxValue)
	greaterSigns := 1
	if equalSigns == progressBarLength {
		greaterSigns = 0
	}
	spaces := progressBarLength - equalSigns - greaterSigns

	return fmt.Sprintf("[%s%s%s] %s/%s",
		strings.Repeat("=", equalSigns),
		strings.Repeat(">", greaterSigns),
		strings.Repeat(" ", spaces),
		cpb.valuePrinter(cpb.currentValue),
		cpb.maxValueString)
}
