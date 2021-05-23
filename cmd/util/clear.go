package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	//dirs := []string{
	//	"4fd66b4b-b976-11eb-816b-e4e74940035b",
	//	"89432f72-b5e9-11eb-a342-e4e74940035b",
	//	"c6c04ca0-b943-11eb-a79d-e4e74940035b",
	//	"1527d749-b86a-11eb-a0d6-e4e74940035b",
	//	"4bf80026-b790-11eb-a175-e4e74940035b",
	//	"de1fa90f-b946-11eb-a79d-e4e74940035b",
	//	"a76d30a5-b866-11eb-a0d6-e4e74940035b",
	//	"85e4778e-b6f3-11eb-98eb-e4e74940035b",
	//	"a64a26b0-b5c0-11eb-b566-e4e74940035b",
	//	"f3c685ea-b942-11eb-a79d-e4e74940035b",
	//	"168dc196-baf9-11eb-bcb4-e4e74940035b",
	//	"c5b34027-b791-11eb-a175-e4e74940035b",
	//	"1e68c398-b5e9-11eb-a342-e4e74940035b",
	//	"21574449-b5e9-11eb-a342-e4e74940035b",
	//	"93264759-b5e9-11eb-a342-e4e74940035b",
	//	"2879f883-b5e9-11eb-a342-e4e74940035b",
	//	"95f88e04-b5e9-11eb-a342-e4e74940035b",
	//	"1a682a20-b5e9-11eb-a342-e4e74940035b",
	//	"af3eee9f-b903-11eb-b645-e4e74940035b",
	//	"cdc1a0b3-b943-11eb-a79d-e4e74940035b",
	//	"f9648da9-b865-11eb-a0d6-e4e74940035b",
	//	"903ae11d-b5e9-11eb-a342-e4e74940035b",
	//	"8be2d07b-b5e9-11eb-a342-e4e74940035b",
	//	"ac7545ec-b947-11eb-be57-e4e74940035b",
	//	"150769da-bac0-11eb-812d-e4e74940035b",
	//	"fb3c38ff-ba87-11eb-a368-e4e74940035b",
	//	"13be6f11-b78e-11eb-a175-e4e74940035b",
	//	"3379d52f-b911-11eb-8df9-e4e74940035b"}
	var connector database.Connector
	err := connector.Connect()
	if err != nil {
		panic("unable to connect to connector" + err.Error())
	}
	defer connector.Close()
	rows, err := connector.Database.Query("SELECT id_video FROM video WHERE quality = ''")
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
