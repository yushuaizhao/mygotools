package mergeFile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Merge use filepath to walk the directory and use ioutil to read,use os.write to append.
func Merge(path string, file string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fmt.Println(path, len(path))
		fmt.Println(f.Size())
		writeTo(path, file)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func writeTo(readFile, writeFile string) {
	var f *os.File
	var err1 error
	buf, err := ioutil.ReadFile(readFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %s\n", err)
	}
	if checkfileExist(writeFile) {
		f, err1 = os.OpenFile(writeFile, os.O_APPEND|os.O_WRONLY, 0666)

	} else {
		f, err1 = os.Create(writeFile)
	}
	if err1 != nil {
		fmt.Println("write failed")
		panic(err1)
	}
	f.Write(buf)
	defer f.Close()
}

func checkfileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
