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
	"strings"
)

const (
	DATAFILE_PATH = "mytodo"
)

type TodoListWrapper struct {
	Version int
	// Items   []TodoListItem `json:"items"`
	Items []LineEntry
	// Items map[int]TodoListItem `json:"-"`
}

// TodoListItem holds a todo list item.
type TodoListItem struct {
	Line uint64
	Body string
}

func (item TodoListItem) String() string {
	return fmt.Sprintf("%d %s", item.Line, item.Body)
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
			myfd.ProcessLine(line)
		}
		myfd.Sort()
	}
}

func (myfd *TodoListWrapper) ProcessLine(s string) {
	if s == "" {
		return
	}

	entry, err := parseLine(s)
	if err != nil {
		log.Fatal(err)
	}

	myfd.Items = Filter(myfd.Items, func(e LineEntry) bool {
		return e.Number != entry.Number
	})

	if !entry.Empty() {
		myfd.Items = append(myfd.Items, entry)
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
	out := make([]string, 0)
	for _, item := range myfd.Items {
		out = append(out, item.String())
	}
	myString := strings.Join(out, "\n")
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

func (myfd *TodoListWrapper) Sort() {
	sort.Sort(ByLineNumber(myfd.Items))
	// sort.Slice(myfd.Items2, func(i, j int) bool {
	// 	return myfd.Items2[i].Line < myfd.Items2[j].Line
	// })
}

func Filter(vs []LineEntry, f func(LineEntry) bool) []LineEntry {
	vsf := make([]LineEntry, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (myfd *TodoListWrapper) ShowItems() {
	for _, item := range myfd.Items {
		fmt.Println(item.String())
	}
}

func main() {
	mydatapath := getDataFilePath()
	dt := TodoListWrapper{}
	dt.readText(mydatapath)

	if len(os.Args) == 2 {
		dt.ProcessLine(os.Args[1])
		dt.Sort()
	}
	dt.ShowItems()
	dt.writeText(mydatapath)
}
