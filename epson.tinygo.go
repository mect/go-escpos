//go:build tinygo

// This file is specifically for making things compatible with TinyGo, it will not run in normal Go!

package escpos

import "io"

// our conversion table does not fit in microcontrollers, so we fake the conversion
// we should find a solution someday
type tinyCharacterConverter struct{}

func (t tinyCharacterConverter) Encode(utf_8 []byte) (latin []byte, success int, err error) {
	return utf_8, 1, nil
}

var converter characterConverter = tinyCharacterConverter{}

type UART interface {
	io.Reader
	io.Writer

	Buffered() int
}

// convert our UART to also have a closer, which it does not have
type uartToRWC struct {
	u UART
}

func (u *uartToRWC) Read(p []byte) (int, error) {
	return u.u.Read(p)
}

func (u *uartToRWC) Write(p []byte) (int, error) {
	return u.u.Write(p)
}

func (u *uartToRWC) Close() error { // fake this out
	return nil
}

// NewPrinterByUART returns a new printer with a TinyGo UART interface
func NewPrinterByUART(uart UART) (*Printer, error) {
	return &Printer{
		s: &uartToRWC{
			u: uart,
		},
	}, nil
}

func (p *Printer) write(cmd string) error {
	_, err := p.s.Write([]byte(cmd))
	return err
}
