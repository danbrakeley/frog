package frog

// ImplData is additional data required to pass from child to parent in order for advanced features
// like anchoirng to function properly.
type ImplData struct {
	// AnchoredLine should be 0 to indicate no anchor, or any number > 0 to uniquely identify a given
	// anchored line
	AnchoredLine int32

	// MinLevel is passed up to the RootLogger, where it is used to decide if this message should be
	// processed or not.
	// Each child that is passed this MinLevel should update it (e.g. via MergeMinLevel) to be the
	// max of the passed MinLevel and its own internal MinLevel.
	MinLevel Level

	// Fields holds Fielders that have already been turned into Fields. This is used by
	// CustomizerLoggers to cache the fields that will be included with every log message.
	Fields []Field
}

// MergeMinLevel sets MinLevel to the max of its own MinLevel and the passed in Level.
func (d *ImplData) MergeMinLevel(min Level) {
	if d.MinLevel < min {
		d.MinLevel = min
	}
}

// MergeFields adds any passed in fields before the existing fields
func (d *ImplData) MergeFields(fields []Field) {
	if len(fields) == 0 {
		return
	}
	d.Fields = append(fields, d.Fields...)
}
