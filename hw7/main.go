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
		// enviroment variable name is the name of file
		varname := f.Name()
		var openerr error
		var file *os.File
		if file, openerr = os.Open(f.Name()); openerr != nil {
			fmt.Println("File", varname, "can't read,", openerr.Error())
			continue
		}
		var stat os.FileInfo
		if stat, openerr = file.Stat(); openerr != nil {
			fmt.Println("File", varname, "can't read,", openerr.Error())
			continue
		}
		if stat.IsDir() {
			continue
		}
		// 32760 - maximum length of enviroment variable
		const maxlen = 32760
		if stat.Size() >= maxlen {
			fmt.Println("Enviroment variable", varname, "can't set. Value is bigger then maximum len", maxlen)
			continue
		}

		file.Close()
	}

	//dd

	fmt.Println(envdir, apppath)
}
