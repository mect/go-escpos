package main

import (
	"fmt"

	"github.com/mect/go-escpos"
)

func main() {
	p, err := escpos.NewUSBPrinter(0, 0)
	if err != nil {
		fmt.Print(err)
		return
	}

	p.Init()
	p.Smooth(true)
	p.Size(2, 2)
	p.PrintLn("HELLO GO")
	p.Size(1, 1)

	p.Font(escpos.FontB)
	p.PrintLn("This is a test of MECT go-escpos")
	p.Font(escpos.FontA)

	p.Align(escpos.AlignRight)
	p.PrintLn("An all Go\neasy to use\nEpson POS Printer library")
	p.Align(escpos.AlignLeft)

	p.Size(2, 2)
	p.PrintLn("* No magic numbers")
	p.PrintLn("* ISO8859-15 ŠÙþþØrt")
	p.Underline(true)
	p.PrintLn("* Extended layout")
	p.Underline(false)
	p.PrintLn("* All in Go!")

	p.Align(escpos.AlignCenter)
	p.Barcode("MECT", escpos.BarcodeTypeCODE39)
	p.Align(escpos.AlignLeft)

	p.Feed(2)
	p.Cut()
	p.End()
}
