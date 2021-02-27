package escpos

// Font defined the font type
type Font uint8

const (
	FontA Font = 0
	FontB Font = 1
	FontC Font = 2
)

// Alignment defines the text blignment
type Alignment uint8

const (
	AlignLeft   Alignment = 0
	AlignCenter Alignment = 1
	AlignRight  Alignment = 2
)

// BarcodeType defines the type of barcode
type BarcodeType string

const (
	// function type A
	BarcodeTypeUPCA    BarcodeType = "\x00"
	BarcodeTypeUPCE    BarcodeType = "\x01"
	BarcodeTypeEAN13   BarcodeType = "\x02"
	BarcodeTypeEAN8    BarcodeType = "\x03"
	BarcodeTypeCODE39  BarcodeType = "\x04"
	BarcodeTypeITF     BarcodeType = "\x05"
	BarcodeTypeCODABAR BarcodeType = "\x06"

	// function type B
	BarcodeTypeCODE128 BarcodeType = "\x49"
)
