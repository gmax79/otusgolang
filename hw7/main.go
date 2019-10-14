package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	envdir := os.Args[1]
	apppath := os.Args[2]
	vars, err := loadEnvVariables(envdir)
	if err != nil {
		log.Fatal(err)
	}
	err = runApp(apppath, os.Args[3:], vars, os.Stdout, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
}

func runApp(app string, args []string, env map[string]string, stdout io.Writer, stderr io.Writer) error {
	for n, v := range env {
		err := os.Setenv(n, v)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command(app, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func loadEnvVariables(envdir string) (map[string]string, error) {
	// read enviroment variables from files in envdir
	envvars := make(map[string]string)
	files, err := ioutil.ReadDir(envdir)
	if err != nil {
		return envvars, err
	}
	for _, f := range files {
		// enviroment variable name is the name of file
		filename := path.Join(envdir, f.Name())
		varname := f.Name()
		var openerr error
		var file *os.File
		if file, openerr = os.Open(filename); openerr != nil {
			return envvars, fmt.Errorf("Error, %s", openerr.Error())
		}
		var stat os.FileInfo
		if stat, openerr = file.Stat(); openerr != nil {
			return envvars, fmt.Errorf("File %s can't read, %s", filename, openerr.Error())
		}
		if stat.IsDir() {
			fmt.Println("dir", filename)
			continue
		}
		const maxlen = 32760 // - maximum length of enviroment variable in os
		if stat.Size() >= maxlen {
			return envvars, fmt.Errorf("Enviroment variable %s can't set. Value is bigger then possible maximum (%d) len", varname, maxlen)
		}
		data, err := ioutil.ReadAll(file)
		file.Close()
		if err != nil {
			return envvars, fmt.Errorf("File %s can't read, %s", filename, openerr.Error())
		}
		// check data correctness
		for _, x := range data {
			if x >= 32 && x <= 127 {
				// ok, do nothing
			} else {
				return envvars, fmt.Errorf("Enviroment variable %s can't set. Value contains invalid symbols", varname)
			}
		}
		envvars[varname] = string(data)
	}
	return envvars, nil
}
