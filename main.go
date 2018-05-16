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

type TodoListWrapper struct {
	Items2 map[uint64]string
}

func NewTodoListWrapper() *TodoListWrapper {
	todoListWrapper := TodoListWrapper{}
	todoListWrapper.Items2 = make(map[uint64]string)
	return &todoListWrapper
}

func (myfd *TodoListWrapper) keys() []uint64 {
	out := make([]uint64, 0)
	for k, _ := range myfd.Items2 {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

func (myfd *TodoListWrapper) String() string {
	out := make([]string, 0)
	for _, k := range myfd.keys() {
		out = append(out, myfd.Items2[k])
	}
	return strings.Join(out, "\n")
}

func (myfd *TodoListWrapper) Add(s string) {
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
	if s != lineNumberString {
		myfd.Items2[lineNumber] = s
	} else {
		delete(myfd.Items2, lineNumber)
	}
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
func (myfd *TodoListWrapper) readText(filepath string) {
	dat := readString(filepath)
	if dat != "" {
		for _, line := range strings.Split(dat, "\n") {
			myfd.Add(line)
		}
	}
}

// ReadJSON reads JSON data from a file.
// func (myfd *TodoListWrapper) readJSON(filepath string) {
// 	dat := readData(filepath)
// 	if dat != nil {
// 		data := []string{}
// 		err := json.Unmarshal(dat, &data)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		myfd.FromArray(data)
// 	}
// }

// Reset clears todo list items.
// func (myfd *TodoListWrapper) Reset() {
// 	myfd.Items = make([]TodoListItem, 0)
// }

// WriteText writes text data to a file.
func (myfd *TodoListWrapper) writeText(filepath string) {
	myString := myfd.String()
	writeString(filepath, myString, 0644)
}

// WriteJSON writes JSON data to a file.
// func (myfd *TodoListWrapper) writeJSON(filepath string) {
// 	strings := myfd.ToArray()
// 	bytes, err := json.Marshal(strings)
// 	if err != nil {
// 		// TODO: try to recover or write temp file here
// 		log.Fatal(err)
// 	}
// 	writeData(filepath, bytes, 0644)
// }

func main() {
	mydatapath := getDataFilePath()
	dt := NewTodoListWrapper()
	dt.readText(mydatapath)

	if len(os.Args) == 2 {
		dt.Add(os.Args[1])
	}
	fmt.Println(dt.String())
	dt.writeText(mydatapath)
}
