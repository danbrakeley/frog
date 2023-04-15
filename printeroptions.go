package frog

type PrinterOption interface {
	isPrinterOption()
	String() string
}

// Palette

func POPalette(p Palette) PrinterOption {
	return poPalette{Palette: p}
}

type poPalette struct {
	Palette Palette
}

func (_ poPalette) isPrinterOption() {}
func (_ poPalette) String() string   { return "POPalette" }

// Time

func POTime(visible bool) poTime {
	return poTime{Visible: visible}
}

type poTime struct {
	Visible bool
}

func (p poTime) isPrinterOption() {}
func (p poTime) String() string   { return "POTime" }

// Level

func POLevel(visible bool) poLevel {
	return poLevel{Visible: visible}
}

type poLevel struct {
	Visible bool
}

func (p poLevel) isPrinterOption() {}
func (p poLevel) String() string   { return "POLevel" }

// Field Indent

func POFieldIndent(indent int) poFieldIndent {
	return poFieldIndent{Indent: indent}
}

type poFieldIndent struct {
	Indent int
}

func (p poFieldIndent) isPrinterOption() {}
func (p poFieldIndent) String() string   { return "POFieldIndent" }

// Message Order: MsgLeftFieldsRight and FieldsLeftMsgRight

var POMsgLeftFieldsRight poMsgLeftFieldsRight

type poMsgLeftFieldsRight struct{}

func (p poMsgLeftFieldsRight) isPrinterOption() {}
func (p poMsgLeftFieldsRight) String() string   { return "POMsgLeftFieldsRight" }

var POFieldsLeftMsgRight poFieldsLeftMsgRight

type poFieldsLeftMsgRight struct{}

func (p poFieldsLeftMsgRight) isPrinterOption() {}
func (p poFieldsLeftMsgRight) String() string   { return "POFieldsLeftMsgRight" }

// Transient Line Length (meant to crop anchored lines so they don't wrap)

func POTransientLineLength(cols int) poTransientLineLength {
	return poTransientLineLength{Cols: cols}
}

type poTransientLineLength struct {
	Cols int
}

func (p poTransientLineLength) isPrinterOption() {}
func (p poTransientLineLength) String() string   { return "POTransientLineLength" }
