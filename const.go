package escpos

// Font defined the font type
type Font uint8

const (
	FontA = 0
	FontB = 1
	FontC = 2
)

// Alignment defines the text blignment
type Alignment uint8

const (
	AlignLeft   = 0
	AlignCenter = 1
	AlignRight  = 2
)

// BarcodeType defines the type of barcode
type BarcodeType string

const (
	// function type A
	BarcodeTypeUPCA    = "\x00"
	BarcodeTypeUPCE    = "\x01"
	BarcodeTypeEAN13   = "\x02"
	BarcodeTypeEAN8    = "\x03"
	BarcodeTypeCODE39  = "\x04"
	BarcodeTypeITF     = "\x05"
	BarcodeTypeCODABAR = "\x06"

	// function type B
	BarcodeTypeCODE128 = "\x49"
)
