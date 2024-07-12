//go:build !tinygo

// This file contains all code that is not compatible with TinyGo

package escpos

import (
	"fmt"
	"image"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bjarneh/latinx"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/aztec"
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

// AztecViaImage prints an Aztec code using the image system for longer data that is not possible to print directly
func (p *Printer) AztecViaImage(data string, width, height int) error {
	if height < 1 {
		height = 500
	}
	if width < 1 {
		width = 500
	}
	aztecCode, err := aztec.Encode([]byte(data), aztec.DEFAULT_EC_PERCENT, aztec.DEFAULT_LAYERS)
	if err != nil {
		return fmt.Errorf("failed to encode aztec code: %w", err)
	}

	// Scale the barcode to 200x200 pixels
	aztecCode, err = barcode.Scale(aztecCode, width, height)
	if err != nil {
		return fmt.Errorf("failed to scale aztec code: %w", err)
	}

	return p.Image(aztecCode)
}

// Image prints a raster image
//
// The image must be narrower than the printer's pixel width
func (p *Printer) Image(img image.Image) error {
	xL, xH, yL, yH, imgData := printImage(img)
	return p.write("\x1dv\x30\x00" + string(append([]byte{xL, xH, yL, yH}, imgData...)))
}
