//go:build !tinygo

// This file contains all code that is not compatible with TinyGo

package escpos

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bjarneh/latinx"
)

var converter characterConverter = latinx.Get(latinx.ISO_8859_15)

// NewUSBPrinter returns a new printer with a USB Vendor and Product ID
// if both are 0 it will return the first found Epson POS printer
func NewUSBPrinterByPath(devpath string) (*Printer, error) {
	if devpath == "" {
		entries, err := os.ReadDir("/dev/usb")
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "lp") {
				devpath = path.Join("/dev/usb", entry.Name())
				break
			}
		}

		if devpath == "" {
			return nil, ErrorNoDevicesFound
		}
	}

	f, err := os.OpenFile(devpath, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("couldn't open %q device: %w", devpath, err)
	}
	return &Printer{
		s: f,
		f: f,
	}, nil
}

func (p *Printer) write(cmd string) error {
	if p.f != nil {
		p.f.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}
	_, err := p.s.Write([]byte(cmd))
	return err
}
