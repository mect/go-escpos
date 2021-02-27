package escpos

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bjarneh/latinx"

	"github.com/mect/go-escpos/pkg/udev"
)

var ErrorNoDevicesFound = errors.New("No devices found")

type Printer struct {
	s io.ReadWriteCloser
}

// NewUSBPrinter returns a new printer with a USb Vendor and Product ID
// if both are 0 it will return the first found Epson POS printer
func NewUSBPrinter(vendorID uint16, productID uint16) (*Printer, error) {
	sc := udev.NewScanner()
	devices, err := sc.ScanDevices()
	if err != nil {
		return nil, fmt.Errorf("couldn't list USB devices: %w", err)
	}

	devpath := ""
	for _, dev := range devices {
		matches := true
		if vendorID != 0 && dev.Env["ID_VENDOR_FROM_DATABASE"] != "" {
			if matches {
				i, _ := strconv.ParseInt(dev.Env["ID_VENDOR_FROM_DATABASE"], 16, 32)
				matches = uint16(i) == vendorID
			}
		}
		if productID != 0 && dev.Env["ID_MODEL_ID"] != "" {
			if matches {
				i, _ := strconv.ParseInt(dev.Env["ID_MODEL_ID"], 16, 32)
				matches = uint16(i) == vendorID
			}
		}

		if vendorID == 0 && productID == 0 && strings.Contains(dev.Env["ID_VENDOR_FROM_DATABASE"], "Epson") {
			matches = true
		}

		if matches {
			if dev.Env["DEVNAME"] != "" {
				devpath = path.Join("/dev/", dev.Env["DEVNAME"])
				break
			}
		}
	}

	if devpath == "" {
		return nil, ErrorNoDevicesFound
	}

	f, err := os.OpenFile("/dev/usb/lp4", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("couldn't open %q device: %w", devpath, err)
	}
	return &Printer{
		s: f,
	}, nil
}

func (p *Printer) write(cmd string) error {
	_, err := p.s.Write([]byte(cmd))
	return err
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
	converter := latinx.Get(latinx.ISO_8859_15)
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

	return p.Print(fmt.Sprintf("%s", barcode))
}
