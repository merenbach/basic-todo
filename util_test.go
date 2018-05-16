package main

import (
	"testing"
)

func TestStringSplitCarCdr(t *testing.T) {
	tables := [][]string{
		{"hello world", " ", "hello", "world"},
		{"hello, world", ", ", "hello", "world"},
		{"hello, world", ",", "hello", " world"},
		{"Hello, world!", " ", "Hello,", "world!"},
	}
	for _, table := range tables {
		s, delim, expected_car, expected_cdr := table[0], table[1], table[2], table[3]
		car, cdr := stringSplitCarCdr(s, delim)
		if car != expected_car || cdr != expected_cdr {
			t.Errorf("Splitting %q on delimiter %q, expected CAR %q (got %q) and CDR %q (got %q)", s, delim, expected_car, car, expected_cdr, cdr)
		}

	}
}

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
