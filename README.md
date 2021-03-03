# go-escpos
This is a Go library to talk to escpos capable devices, mainly the Epson POS printers. 

it will allow you to send print jobs in a simple functional interface, this library also supports automatic discovery of USB devices.

[Documentation can be found on pkg.go.dev](https://pkg.go.dev/github.com/mect/go-escpos)

## Example

```go
package main

import (
	"fmt"

	"github.com/mect/go-escpos"
)

func main() {
	p, err := escpos.NewUSBPrinterByPath("") // empry string will do a self discovery
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Init() // start 
	p.Smooth(true) // use smootth printing
	p.Size(2, 2) // set font size
	p.PrintLn("HELLO GO")

	p.Size(1, 1)
	p.Font(escpos.FontB) // change font
	p.PrintLn("This is a test of MECT go-escpos")
	p.Font(escpos.FontA)

	p.Align(escpos.AlignRight) // change alignment
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
	p.Barcode("MECT", escpos.BarcodeTypeCODE39) // print barcode
	p.Align(escpos.AlignLeft)

	p.FeedN(2) // feed 2
	p.Cut() // cut
	p.End() // stop
}

```

## Limitations
* It currently only supports USB
* The USB discovery only works on Linux due to the use of `udev`

(PRs to fix these welcome!)

## Known to work with
* Epson TM-20III