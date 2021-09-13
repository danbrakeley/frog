package frog

type PrinterOption interface {
	Value() PO
	AsBool() bool
	AsInt() int
}

type PO int

func (_ PO) isPrinterOption() {}

const (
	poPalette PO = iota
	poShowTime
	poShowLevel
	poFieldIndent
	poMessageLast
)

// Palette

func POPalette(pal Palette) PrinterOption {
	return potPalette(pal)
}

type potPalette Palette

func (p potPalette) Value() PO {
	return poPalette
}

func (p potPalette) AsBool() bool {
	return false
}

func (p potPalette) AsInt() int {
	return int(p)
}

// Show Time

func POShowTime(show bool) PrinterOption {
	return potShowTime(show)
}

type potShowTime bool

func (p potShowTime) Value() PO {
	return poShowTime
}

func (p potShowTime) AsBool() bool {
	return bool(p)
}

func (p potShowTime) AsInt() int {
	return 0
}

// Show Level

func POShowLevel(show bool) PrinterOption {
	return potShowLevel(show)
}

type potShowLevel bool

func (p potShowLevel) Value() PO {
	return poShowLevel
}

func (p potShowLevel) AsBool() bool {
	return bool(p)
}

func (p potShowLevel) AsInt() int {
	return 0
}

// Field Indent

func POFieldIndent(indent int) PrinterOption {
	return potFieldIndent(indent)
}

type potFieldIndent int

func (p potFieldIndent) Value() PO {
	return poFieldIndent
}

func (p potFieldIndent) AsBool() bool {
	return false
}

func (p potFieldIndent) AsInt() int {
	return int(p)
}

// Field Indent

func POMessageLast(b bool) PrinterOption {
	return potMessageLast(b)
}

type potMessageLast bool

func (p potMessageLast) Value() PO {
	return poMessageLast
}

func (p potMessageLast) AsBool() bool {
	return bool(p)
}

func (p potMessageLast) AsInt() int {
	return 0
}
