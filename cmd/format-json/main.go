package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

func die(err error) {
	panic(err.Error())
}

func format(inFile string, inPlace bool) error {
	bs, err := ioutil.ReadFile(inFile)
	if err != nil {
		return (err)
	}

	var obj interface{}
	if err := yaml.Unmarshal(bs, &obj); err != nil {
		return (err)
	}

	switch filepath.Ext(inFile) {
	case ".yaml":
		fallthrough
	case ".yml":
		bs, err = yaml.Marshal(obj)
	case ".json":
		fallthrough
	default:
		bs, err = json.MarshalIndent(obj, "", "  ")
	}
	if err != nil {
		return err
	}

	if !inPlace {
		io.Copy(os.Stdout, bytes.NewBuffer(bs))
		return nil
	}

	f, err := ioutil.TempFile(".", "format")
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, bytes.NewBuffer(bs)); err != nil {
		f.Close()
		return err
	}
	f.Close()
	if err := os.Rename(f.Name(), inFile); err != nil {
		return err
	}
	return nil
}

func main() {
	var inPlace bool
	flag.BoolVar(&inPlace, "inPlace", false, "update file in place")

	flag.Parse()

	for _, arg := range flag.Args() {
		fi, err := os.Stat(arg)
		if err != nil {
			die(err)
		}

		if !fi.IsDir() {
			if err := format(arg, inPlace); err != nil {
				die(fmt.Errorf("error formatting %q: %s", arg, err))
			}
		} else {
			if err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				switch filepath.Ext(path) {
				case ".yaml":
					fallthrough
				case ".yml":
					fallthrough
				case ".json":
					if err := format(path, inPlace); err != nil {
						return fmt.Errorf("error formatting %q: %s", path, err)
					}
					return nil
				}
				return nil
			}); err != nil {
				die(err)
			}
		}
	}
}
