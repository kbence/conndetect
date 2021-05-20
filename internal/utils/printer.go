package utils

import "fmt"

type Printer interface {
	Printf(format string, args ...interface{})
}

type PrinterImpl struct{}

func NewPrinter() Printer {
	return &PrinterImpl{}
}

func (p *PrinterImpl) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
