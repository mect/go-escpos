package escpos

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrorNoDevicesFound = errors.New("No devices found")

// fufilled by either tinygoConverter or latinx
type characterConverter interface {
	Encode(utf_8 []byte) (latin []byte, success int, err error)
}

type Printer struct {
	s io.ReadWriteCloser
	f *os.File
}

func NewPrinterByRW(rwc io.ReadWriteCloser) (*Printer, error) {

	return &Printer{
		s: rwc,
	}, nil
}

// Init sends an init signal
func (p *Printer) Init() error {
	// send init command
	err := p.write("\x1B@")
	if err != nil {
		return err
	}

	// send encoding ISO8859-15
	return p.write(fmt.Sprintf("\x1Bt%c", 40))
}

// End sends an end signal to finalize the print job
func (p *Printer) End() error {
	return p.write("\xFA")
}

// Close closes the connection to the printer, all commands will not work after this
func (p *Printer) Close() error {
	return p.s.Close()
}

// Cut sends the command to cut the paper
func (p *Printer) Cut() error {
	return p.write("\x1DVA0")
}

// Feed sends a paper feed command for a specified length
func (p *Printer) Feed(n int) error {
	return p.write(fmt.Sprintf("\x1Bd%c", n))
}

// Print prints a string
// the data is re-encoded from Go's UTF-8 to ISO8859-15
func (p *Printer) Print(data string) error {
	if data == "" {
		return nil
	}

	b, _, err := converter.Encode([]byte(data))
	if err != nil {
		return err
	}
	data = string(b)

	data = textReplace(data)

	return p.write(data)
}

// PrintLn does a Print with a newline attached
func (p *Printer) PrintLn(data string) error {
	err := p.Print(data)
	if err != nil {
		return err
	}

	return p.write("\n")
}

// Size changes the font size
func (p *Printer) Size(width, height uint8) error {
	// sended size is 8 bit, 4 width + 4 height
	return p.write(fmt.Sprintf("\x1D!%c", ((width-1)<<4)|(height-1)))
}

// Font changest the font face
func (p *Printer) Font(font Font) error {
	return p.write(fmt.Sprintf("\x1BM%c", font))
}

// Underline will enable or disable underlined text
func (p *Printer) Underline(enabled bool) error {
	if enabled {
		return p.write(fmt.Sprintf("\x1B-%c", 1))
	}
	return p.write(fmt.Sprintf("\x1B-%c", 0))
}

// Smooth will enable or disable smooth text printing
func (p *Printer) Smooth(enabled bool) error {
	if enabled {
		return p.write(fmt.Sprintf("\x1Db%c", 1))
	}
	return p.write(fmt.Sprintf("\x1Db%c", 0))
}

// Align will change the text alignment
func (p *Printer) Align(align Alignment) error {
	return p.write(fmt.Sprintf("\x1Ba%c", align))
}

// Barcode will print a barcode of a specified type as well as the text value
func (p *Printer) Barcode(barcode string, format BarcodeType) error {

	// set width/height to default
	err := p.write("\x1d\x77\x04\x1d\x68\x64")
	if err != nil {
		return err
	}

	// set barcode font
	err = p.write("\x1d\x66\x00")
	if err != nil {
		return err
	}

	switch format {
	case BarcodeTypeUPCA:
		fallthrough
	case BarcodeTypeUPCE:
		fallthrough
	case BarcodeTypeEAN13:
		fallthrough
	case BarcodeTypeEAN8:
		fallthrough
	case BarcodeTypeCODE39:
		fallthrough
	case BarcodeTypeITF:
		fallthrough
	case BarcodeTypeCODABAR:
		err = p.write(fmt.Sprintf("\x1d\x6b%s%v\x00", format, barcode))
	case BarcodeTypeCODE128:
		err = p.write(fmt.Sprintf("\x1d\x6b%s%v%v\x00", format, len(barcode), barcode))
	default:
		panic("unimplemented barcode")
	}

	if err != nil {
		return err
	}

	return p.PrintLn(fmt.Sprintf("%s", barcode))
}

func (p *Printer) GetErrorStatus() (ErrorStatus, error) {
	_, err := p.s.Write([]byte{0x10, 0x04, 0x02})
	if err != nil {
		return 0, err
	}
	data := make([]byte, 1)
	_, err = p.s.Read(data)
	if err != nil {
		return 0, err
	}

	return ErrorStatus(data[0]), nil
}
