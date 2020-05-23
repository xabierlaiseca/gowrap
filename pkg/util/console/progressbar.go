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
	cpb := consoleProgressBar{
		maxValue:        max,
		currentValue:    0,
		maxProgress:     defaultProgressBarLength,
		currentProgress: 0,
		valuePrinter:    printer,
		maxValueString:  printer(max),
	}

	fmt.Print(cpb.buildConsoleLine())
	return &cpb
}

type consoleProgressBar struct {
	maxValue        int64
	currentValue    int64
	maxProgress     int
	currentProgress int
	valuePrinter    Printer

	maxValueString string
}

func (cpb *consoleProgressBar) Increment(inc int64) {
	if cpb.currentValue+inc > cpb.maxValue {
		cpb.currentValue = cpb.maxValue
	} else {
		cpb.currentValue += inc
	}

	previousProgress := cpb.currentProgress
	cpb.currentProgress = int(int64(cpb.maxProgress) * cpb.currentValue / cpb.maxValue)

	if cpb.currentProgress > previousProgress {
		line := cpb.buildConsoleLine()
		fmt.Printf("\r\033[K%s", line)
	}
}

func (cpb *consoleProgressBar) Done() {
	if cpb.currentValue < cpb.maxValue {
		cpb.Increment(cpb.maxValue - cpb.currentValue)
	}

	fmt.Println()
}

func (cpb *consoleProgressBar) buildConsoleLine() string {
	greaterSigns := 1
	if cpb.currentProgress == cpb.maxProgress {
		greaterSigns = 0
	}
	spaces := cpb.maxProgress - cpb.currentProgress - greaterSigns

	return fmt.Sprintf("[%s%s%s] %s/%s",
		strings.Repeat("=", cpb.currentProgress),
		strings.Repeat(">", greaterSigns),
		strings.Repeat(" ", spaces),
		cpb.valuePrinter(cpb.currentValue),
		cpb.maxValueString)
}
