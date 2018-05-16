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

func (entry LineEntry) Empty() bool {
	return entry.Contents == ""
}

func (entry LineEntry) String() string {
	return fmt.Sprintf("%d %s", entry.Number, entry.Contents)
}

// StringSplitCarCdr splits a string into a head (one group of characters) and tail (remaining characters).
func stringSplitCarCdr(s, delim string) (string, string) {
	components := strings.SplitN(s, delim, 2)
	return components[0], strings.Join(components[1:], "")
}

// ParseLine parses a line into a line number prefix and string suffix.
func parseLine(line string) (entry LineEntry, err error) {
	car, cdr := stringSplitCarCdr(strings.TrimSpace(line), " ")

	lineNumber, err := strconv.ParseUint(car, 10, 64)
	if err != nil {
		return entry, err
	}

	entry.Number = lineNumber
	if cdr != "" {
		entry.Contents = strings.TrimSpace(cdr)
	}

	return entry, nil
}

// type LineNumShell map[int]string

// func (lns *LineNumShell) OrderedKeys() []string {
// 	keys := reflect.ValueOf(lns).MapKeys()
// 	k2 := keys[0]
// 	return sort.Strings(k2)
// }
