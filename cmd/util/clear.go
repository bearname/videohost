package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var connector mysql.ConnectorImpl
	err := connector.Connect()
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
	}
	defer connector.Close()
	rows, err := connector.database.Query("SELECT id_video FROM video WHERE quality = ''")
	videoDirs := make([]string, 0)
	for rows.Next() {
		var videoId string
		err := rows.Scan(&videoId)

		if err == nil {
			videoDirs = append(videoDirs, videoId)
		}
	}

	root := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content\\"

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(videoDirs)
	fmt.Println(len(files))
	for _, f := range files {
		if !contains(videoDirs, f.Name()) {
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
