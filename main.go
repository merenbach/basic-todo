package main

// TODO: allow datafile override
// TODO: allow DynamoDB remote storage
// TODO: allow Web app (daemon) mode, incl. w/DynamoDB storage
// TODO: use pointers for memory efficiency
// TODO: use interface so that you can have a local file todo list item, a Dynamo todo list item, etc.
// BASIC todo: enter with line number (or have autogen); ./mytodo list lists; ./mytodo 42 deletes item 42, ./mytodo 42 do stuff replaces or inserts item 42 (with feedback)

import (
	"container/list"
	"encoding/json"
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
	Items  list.List `json:"items"`
	Items2 []TodoListItem
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

// WriteData writes data
func writeData(filepath string, data []byte, mode os.FileMode) {
	err := ioutil.WriteFile(filepath, data, mode)
	if err != nil {
		log.Fatal(err)
	}
}

// ReadText reads text data from a file.
func (myfd *TodoListWrapper) readText(filepath string) {
	dat := readData(filepath)
	if dat != nil {
		for _, line := range strings.Split(string(dat), "\n") {
			myfd.AddItem(line)
		}
	}
}

// ReadJSON reads JSON data from a file.
func (myfd *TodoListWrapper) readJSON(filepath string) {
	dat := readData(filepath)
	if dat != nil {
		data := []string{}
		err := json.Unmarshal(dat, &data)
		if err != nil {
			log.Fatal(err)
		}
		myfd.FromArray(data)
	}
}

// Reset clears todo list items.
// func (myfd *TodoListWrapper) Reset() {
// 	myfd.Items = make([]TodoListItem, 0)
// }

// WriteText writes text data to a file.
func (myfd *TodoListWrapper) writeText(filepath string) {
	s := myfd.ToArray()
	myString := strings.Join(s, "\n")
	writeData(filepath, []byte(myString), 0644)
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

// // ByAge implements sort.Interface for []Person based on
// // the Age field.
// type ByLine []TodoListItem

// func (a ByLine) Len() int           { return len(a) }
// func (a ByLine) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a ByLine) Less(i, j int) bool { return a[i].Line < a[j].Line }
// sort.Sort(ByLine(myfd.Items2))

func (myfd *TodoListWrapper) AddItem(lineEntry LineEntry) {
	// tli := TodoListItem{Body: item}
	// myfd.Items = append(myfd.Items, tli)

	myfd.DelItem(lineEntry.Number)
	myfd.Sort()
	if len(out) < 2 {
		return

	}

	item := TodoListItem{
		Line: linenum,
		Body: out[1],
	}
	myfd.Items2 = append(myfd.Items2, item)
	// fmt.Printf("Inserting %v\n", item)

	// fmt.Println("items now equals: ", myfd.Items2)
	myfd.Sort()
}

func (myfd *TodoListWrapper) Sort() {
	sort.Slice(myfd.Items2, func(i, j int) bool {
		return myfd.Items2[i].Line < myfd.Items2[j].Line
	})
}

func Filter(vs []TodoListItem, f func(TodoListItem) bool) []TodoListItem {
	vsf := make([]TodoListItem, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (myfd *TodoListWrapper) DelItem(lineNum uint64) {

	myfd.Items2 = Filter(myfd.Items2, func(tli TodoListItem) bool {
		return tli.Line != lineNumber
	})

	// fmt.Println("items now equals: ", myfd.Items2)

	// tli := TodoListItem{Body: item}
	// myfd.Items = append(myfd.Items, tli)
	// myfd.Items.Remove(item)
}

func (myfd *TodoListWrapper) ToArray() []string {
	out := make([]string, 0)
	for _, item := range myfd.Items2 {
		out = append(out, item.String())
	}
	// for temp := myfd.Items.Front(); temp != nil; temp = temp.Next() {
	// 	out = append(out, temp.Value.(string))
	// }
	return out
}

func (myfd *TodoListWrapper) FromArray(ss []string) {
	myfd.Items2 = make([]TodoListItem, 0)
	for _, line := range ss {
		myfd.AddItem(line)
	}
	myfd.Sort()
	// myfd.Items.Init()
	// for _, s := range ss {
	// 	myfd.Items.PushBack(s)
	// }
}

func (myfd *TodoListWrapper) ShowItems() {
	for _, item := range myfd.Items2 {
		fmt.Println(item.String())
	}
	// counter := 0
	// for temp := myfd.Items.Front(); temp != nil; temp = temp.Next() {
	// 	counter++
	// 	fmt.Println(counter, temp.Value)
	// }
}

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatalf("USAGE: %s COMMAND", os.Args[0])
	// }
	mydatapath := getDataFilePath()
	dt := TodoListWrapper{}
	dt.readText(mydatapath)
	// if len(os.Args) == 1 {
	// 	dt.ShowItems()
	// 	os.Exit(0)
	// }
	if len(os.Args) == 2 {
		lineEntry := parseLine(os.Args[1])
		fmt.Println(lineEntry)
		dt.AddItem(lineEntry)
		dt.writeText(mydatapath)
	}
	dt.ShowItems()

	// dt.AddItem("42")
	// dt.AddItem("42 pet the cats")
	// dt.AddItem("38 wash the car")
	// dt.AddItem("57 clean the tub")
	// dt.AddItem("823 wear a shirt")
	// fmt.Println(dt.ToArray())
	// dt.DelItem("823")

	// if len(os.Args) == 2 {
	// 	dt.AddItem(os.Args[1])

	// }

	// sample data
	// dt.Reset()
	// newItem1 := TodoListItem{Body: "Household chores"}
	// newItem1_1 := TodoListItem{Body: "Mow the lawn"}
	// newItem1_2 := TodoListItem{Body: "Wash the car"}
	// newItem1.Items = append(newItem1.Items, newItem1_1)
	// newItem1.Items = append(newItem1.Items, newItem1_2)
	// dt.Items = append(dt.Items, newItem1)
	// dt.Items = list.New()
	// dt.Items.PushBack("Mow the lawn")
	// dt.Items.PushBack("Wash tho car")

}
