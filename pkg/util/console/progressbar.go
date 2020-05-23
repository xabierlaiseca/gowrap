package console

import (
	"fmt"
	"strings"
)

const defaultProgressBarLength = 50

type ProgressBar interface {
	Increment(inc int64)
	Done()
}

type Printer func(int64) string

func NewProgressBar(max int64, printer Printer) ProgressBar {
	cpb := &consoleProgressBar{
		maxValue:         max,
		currentValue:     0,
		maxProgress:      defaultProgressBarLength,
		currentProgress:  0,
		valuePrinter:     printer,
		maxFriendlyValue: printer(max),
	}

	cpb.printProgressBar()
	return cpb
}

type consoleProgressBar struct {
	currentValue int64
	maxValue     int64

	currentProgress int
	printedProgress int
	maxProgress     int

	currentFriendlyValue string
	printedFriendlyValue string
	maxFriendlyValue     string

	valuePrinter Printer
}

func (cpb *consoleProgressBar) Increment(inc int64) {
	if cpb.currentValue+inc > cpb.maxValue {
		cpb.currentValue = cpb.maxValue
	} else {
		cpb.currentValue += inc
	}

	cpb.printProgressBar()
}

func (cpb *consoleProgressBar) Done() {
	if cpb.currentValue < cpb.maxValue {
		cpb.Increment(cpb.maxValue - cpb.currentValue)
	}

	fmt.Println()
}

func (cpb *consoleProgressBar) printProgressBar() {
	cpb.currentProgress = int(int64(cpb.maxProgress) * cpb.currentValue / cpb.maxValue)
	cpb.currentFriendlyValue = cpb.valuePrinter(cpb.currentValue)

	if cpb.needsUpdate() {
		greaterSigns := 1
		if cpb.currentProgress == cpb.maxProgress {
			greaterSigns = 0
		}
		spaces := cpb.maxProgress - cpb.currentProgress - greaterSigns

		line := fmt.Sprintf("\r\033[K[%s%s%s] %s/%s",
			strings.Repeat("=", cpb.currentProgress),
			strings.Repeat(">", greaterSigns),
			strings.Repeat(" ", spaces),
			cpb.currentFriendlyValue,
			cpb.maxFriendlyValue)

		fmt.Print(line)

		cpb.printedProgress = cpb.currentProgress
		cpb.printedFriendlyValue = cpb.currentFriendlyValue
	}
}

func (cpb *consoleProgressBar) needsUpdate() bool {
	return cpb.printedProgress != cpb.currentProgress ||
		strings.Compare(cpb.printedFriendlyValue, cpb.currentFriendlyValue) != 0
}
