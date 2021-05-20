package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	dirs := []string{
		"0699cea9-b91d-11eb-a2fc-e4e74940035b",
		"10bcdc26-b871-11eb-a0d6-e4e74940035b",
		"13be6f11-b78e-11eb-a175-e4e74940035b",
		"150f8a85-b870-11eb-a0d6-e4e74940035b",
		"1527d749-b86a-11eb-a0d6-e4e74940035b",
		"1a682a20-b5e9-11eb-a342-e4e74940035b",
		"1e68c398-b5e9-11eb-a342-e4e74940035b",
		"21574449-b5e9-11eb-a342-e4e74940035b",
		"24599ef7-b5e9-11eb-a342-e4e74940035b",
		"2879f883-b5e9-11eb-a342-e4e74940035b",
		"2c5fab6a-b5e9-11eb-a342-e4e74940035b",
		"3379d52f-b911-11eb-8df9-e4e74940035b",
		"3cb8c447-b5ca-11eb-b22b-e4e74940035b",
		"40ae40c6-b91d-11eb-a2fc-e4e74940035b",
		"41bf7f08-b5d6-11eb-b2c5-e4e74940035b",
		"4bf80026-b790-11eb-a175-e4e74940035b",
		"57655d3e-b92d-11eb-989d-e4e74940035b",
		"64736bf5-b449-11eb-875e-00ff7c2a75d7",
		"7e7ca286-b5c6-11eb-a83a-e4e74940035b",
		"85e4778e-b6f3-11eb-98eb-e4e74940035b",
		"89432f72-b5e9-11eb-a342-e4e74940035b",
		"8be2d07b-b5e9-11eb-a342-e4e74940035b",
		"903ae11d-b5e9-11eb-a342-e4e74940035b",
		"93264759-b5e9-11eb-a342-e4e74940035b",
		"95f88e04-b5e9-11eb-a342-e4e74940035b",
		"9fe60de8-b6d6-11eb-9200-e4e74940035b",
		"a64a26b0-b5c0-11eb-b566-e4e74940035b",
		"a76d30a5-b866-11eb-a0d6-e4e74940035b",
		"af3eee9f-b903-11eb-b645-e4e74940035b",
		"b2b8ba78-b86f-11eb-a0d6-e4e74940035b",
		"b591502a-b91c-11eb-a2fc-e4e74940035b",
		"c5b34027-b791-11eb-a175-e4e74940035b",
		"cf561af0-b5f0-11eb-a7d7-e4e74940035b",
		"e016b2cf-b5d6-11eb-b729-e4e74940035b",
		"f9648da9-b865-11eb-a0d6-e4e74940035b"}
	root := "C:\\Users\\mikha\\go\\src\\videohost\\videoserver\\content\\"

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(dirs)
	fmt.Println(len(files))
	for _, f := range files {
		if !contains(dirs, f.Name()) {
			err := RemoveContents(root + f.Name())
			if err != nil {
				fmt.Println(f.Name())
				fmt.Println(err.Error())
			}
			dir := root + f.Name()
			fmt.Println(dir)
			err = os.Remove(dir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
