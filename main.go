package main

// TODO: allow datafile override
// TODO: allow DynamoDB remote storage
// TODO: allow Web app (daemon) mode, incl. w/DynamoDB storage
// TODO: use pointers for memory efficiency
// TODO: use interface so that you can have a local file todo list item, a Dynamo todo list item, etc.
// TODO: use interface{} for key, instead of string or number? how do we handle sorting properly?
// BASIC todo: enter with line number (or have autogen); ./mytodo list lists; ./mytodo 42 deletes item 42, ./mytodo 42 do stuff replaces or inserts item 42 (with feedback)

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"sort"
	"strconv"
	"strings"
)

const (
	DATAFILE_PATH = "mytodo"
)

// TodoList holds a todo list.
type TodoList struct {
	Items map[uint64]string
}

// NewTodoList returns a newly-initialized todo list.
func NewTodoList() *TodoList {
	t := TodoList{}
	t.Init()
	return &t
}

// Keys returns an ordered list of dictionary keys for the todo list.
func (t *TodoList) keys() []uint64 {
	out := make([]uint64, 0)
	for k, _ := range t.Items {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

func (t *TodoList) String() string {
	out := make([]string, 0)
	for _, k := range t.keys() {
		out = append(out, t.Items[k])
	}
	return strings.Join(out, "\n")
}

// Init initalizes a todo list.
func (t *TodoList) Init() {
	t.Items = make(map[uint64]string)
}

// Set changes an item in the todo list.
func (t *TodoList) Set(k uint64, v string) {
	if v != "" {
		t.Items[k] = v
	} else {
		delete(t.Items, k)
	}
}

// Parse parses a line and adds it to the todo list.
func (t *TodoList) Parse(s string) {
	if s == "" {
		return
	}

	lineNumberString := strings.SplitN(s, " ", 2)[0]
	lineNumber, err := strconv.ParseUint(lineNumberString, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// if s is lineNumberString, only a line number was entered
	// [TODO] improve this
	if s == lineNumberString {
		s = ""
	}
	t.Set(lineNumber, s)
}

// GetCurrentUserHomeDir returns the home directory of the user running the program.
func getCurrentUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

// GetDataFilePath returns the configuration file for the program.
func getDataFilePath() string {
	base := getCurrentUserHomeDir()
	dataPath := path.Join(base, DATAFILE_PATH)
	// err := os.MkdirAll(dataPath, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return dataPath
}

// ReadData reads data from a given filepath and will return `nil` if the file doesn't exist.
func readData(filepath string) []byte {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Fatal(err)
	}
	return dat
}

// WriteString reads textual data into a string for a given filepath.
func writeString(filepath string, data string, mode os.FileMode) {
	writeData(filepath, []byte(data), mode)
}

// ReadString reads textual data into a string for a given filepath.
func readString(filepath string) string {
	if data := readData(filepath); data != nil {
		return string(data)
	}
	return ""
}

// WriteData writes data
func writeData(filepath string, data []byte, mode os.FileMode) {
	err := ioutil.WriteFile(filepath, data, mode)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadText reads text data from a file.
func (t *TodoList) readText(filepath string) {
	dat := readString(filepath)
	if dat != "" {
		for _, line := range strings.Split(dat, "\n") {
			t.Parse(line)
		}
	}
}

// ReadJSON reads JSON data from a file.
// func (t *TodoList) readJSON(filepath string) {
// 	dat := readData(filepath)
// 	if dat != nil {
// 		data := []string{}
// 		err := json.Unmarshal(dat, &data)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		t.FromArray(data)
// 	}
// }

// Reset clears todo list items.
// func (t *TodoList) Reset() {
// 	t.Items = make([]TodoListItem, 0)
// }

// WriteText writes text data to a file.
func (t *TodoList) writeText(filepath string) {
	myString := t.String()
	writeString(filepath, myString, 0644)
}

// WriteJSON writes JSON data to a file.
// func (t *TodoList) writeJSON(filepath string) {
// 	strings := t.ToArray()
// 	bytes, err := json.Marshal(strings)
// 	if err != nil {
// 		// TODO: try to recover or write temp file here
// 		log.Fatal(err)
// 	}
// 	writeData(filepath, bytes, 0644)
// }

func main() {
	mydatapath := getDataFilePath()
	dt := NewTodoList()
	dt.readText(mydatapath)

	if len(os.Args) == 2 {
		dt.Parse(os.Args[1])
	}
	fmt.Println(dt.String())
	dt.writeText(mydatapath)
}
