package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func usage(fullpath string) {
	appname := path.Base(fullpath)
	fmt.Printf(`Run other program with custom enviroment
Usage:
	./%s envdir application [parameters]"
	envdir - directory with text files, where declared enviroment variables
	application - executable file path
	parameters - custom parameters for application
`, appname)
}

func main() {
	count := len(os.Args)
	if count < 3 {
		usage(os.Args[0])
		return
	}

	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	envdir := os.Args[1]
	apppath := os.Args[2]

	// read enviroment variables from files in envdir
	files, err := ioutil.ReadDir(envdir)
	if err != nil {
		return
	}

	for _, f := range files {
		var file *os.File
		var openerr error
		if file, openerr = os.Open(f.Name()); openerr != nil {
			fmt.Println("File", f.Name, "can't read,", openerr.Error())
			continue
		}
		file.Close()

	}

	fmt.Println(envdir, apppath)
}
