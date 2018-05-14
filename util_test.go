package main

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tables := map[string]LineEntry{
		"10 REM This is a test": LineEntry{
			10,
			"REM This is a test",
		},
		"20 INPUT S": LineEntry{
			20,
			"INPUT S",
		},

		"30 PRINT S": LineEntry{
			30,
			"PRINT S",
		},
		"  10 REM Initial space": LineEntry{
			10,
			"REM Initial space",
		},
		"40    END": LineEntry{
			40,
			"END",
		},
		"20  ": LineEntry{
			20,
			"",
		},
		"30 ": LineEntry{
			30,
			"",
		},
		"   40   ": LineEntry{
			40,
			"",
		},
		"40": LineEntry{
			40,
			"",
		},
	}
	for k, v := range tables {
		entry, _ := parseLine(k)
		if entry.Number != v.Number || entry.Contents != v.Contents {
			t.Fatalf("Expected line entry %q but got %q\n", v, entry)
		}
	}
}
