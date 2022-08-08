package main

import (
	"machine"
	"time"

	"github.com/mect/go-escpos"
)

// these are the default pins for the Arduino Nano33 IoT.
var (
	uart = machine.UART2
	tx   = machine.PB22
	rx   = machine.PB23

	console = machine.Serial
)

func main() {
	time.Sleep(time.Second)

	print("HELLO TINY GOPHERS")

	uart.Configure(machine.UARTConfig{
		TX:       tx,
		RX:       rx,
		BaudRate: 38400, // you can get this by holding the feed button when turning on the printer
	})

	p, _ := escpos.NewPrinterByUART(uart)

	p.Init()

	p.Smooth(true)
	p.Size(2, 2)
	p.PrintLn("HELLO TinyGo")
	p.Size(1, 1)

	p.Font(escpos.FontB)
	p.PrintLn("This is a test of MECT/go-escpos")
	p.Font(escpos.FontA)

	p.Align(escpos.AlignRight)
	p.PrintLn("An all Go\neasy to use\nEpson POS Printer library")
	p.Align(escpos.AlignLeft)

	p.Size(2, 2)
	p.PrintLn("* No magic numbers")
	p.PrintLn("* One codebase for Go\n   and TinyGo")
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

	p.End()

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// we're done, let's blink some LEDs while waitng for a reboot
	for {
		led.Low()

		time.Sleep(time.Millisecond * 500)

		led.High()
		time.Sleep(time.Millisecond * 500)
	}

}
