package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByLineNumber []LineEntry

func (a ByLineNumber) Len() int           { return len(a) }
func (a ByLineNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLineNumber) Less(i, j int) bool { return a[i].Number < a[j].Number }

// sort.Sort(ByLineNumber(myfd.Items2))

// LineEntry holds a line number and line contents.
type LineEntry struct {
	Number   uint64
	Contents string
}

func (entry LineEntry) String() string {
	return fmt.Sprintf("%d %s", entry.Number, entry.Contents)
}

// ParseLine parses a line into a line number prefix and string suffix.
func parseLine(line string) (LineEntry, error) {
	entry := LineEntry{}

	components := strings.SplitN(strings.TrimSpace(line), " ", 2)
	lineNumber, err := strconv.ParseUint(components[0], 10, 64)
	if err != nil {
		return entry, err
	}

	entry.Number = lineNumber
	if len(components) == 2 {
		entry.Contents = strings.TrimSpace(components[1])
	}

	return entry, nil
}

// type LineNumShell map[int]string

// func (lns *LineNumShell) OrderedKeys() []string {
// 	keys := reflect.ValueOf(lns).MapKeys()
// 	k2 := keys[0]
// 	return sort.Strings(k2)
// }
