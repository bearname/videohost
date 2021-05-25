package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	root := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content\\"

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(root)
	fmt.Println(len(files))
	for _, f := range files {
		if f.IsDir() {
			videoDir := root + "/" + f.Name()
			filesInDir, err := ioutil.ReadDir(videoDir)
			if err != nil {
				log.Error(err)
			}
			for _, file := range filesInDir {
				name := file.Name()
				oldPath := videoDir + "/" + name
				if name[0] == 'i' && name[1] == 'n' && name[2] == 'd' && name[3] == 'e' && name[4] == 'x' && name[5] != '-' {
					i := len(name)
					s := name[i-len(".mp4") : i]

					newPath := videoDir + "index-1080" + s
					rename(oldPath, newPath)
				}

				if strings.Contains(name, "index-1080x") {
					i := len(name)
					s := name[i-len(".t"): i]

					newPath := videoDir + "/" + "index-1080" + s + ".ts"
					rename(oldPath, newPath)
					fmt.Println(oldPath)
					fmt.Println(newPath)
				}
			}
		}
	}
}

func rename(oldPath string, newPath string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("Failed rename " + oldPath + " to " + newPath)
	}
}
