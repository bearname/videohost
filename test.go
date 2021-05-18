package main

import (
	"fmt"
	"github.com/bearname/videohost/videoserver/util"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage test.exe <browser tab count> <count video server>")
		return
	}
	toInt, b := util.StrToInt(os.Args[1])
	if !b {
		fmt.Println("Invalid browser tab count")
		return
	}
	countVideoServer, ok := util.StrToInt(os.Args[2])
	if !ok {
		fmt.Println("Invalid count video server")
		return
	}
	fmt.Println("countVideoServer " + strconv.Itoa(countVideoServer))
	min := 8000
	max := min + countVideoServer
	//for i := min; i < max; i++ {
	//	fmt.Println("port " + Itoa(i))
	//	port :=Itoa(i)
	//	itoa := Itoa(min)
	//	s := Itoa(max)
	//	fmt.Println(s + " " + itoa)
	//	err := exec.Command("runserver.exe", itoa, s).Run()
	//	var message string
	//	if err != nil {
	//		message = "failed "
	//	} else {
	//		message = "success "
	//	}
	//	fmt.Println(message + " run videoserver on port " + port)
	//}

	for i := 0; i < toInt; i++ {
		rand.Seed(time.Now().UnixNano())
		port := rand.Intn(max-min) + min
		fmt.Println(port)
		url := "http://localhost:" + Itoa(port) + "/videos/cf561af0-b5f0-11eb-a7d7-e4e74940035b"
		fmt.Println(url)
		command := exec.Command("C:\\Program Files\\Mozilla Firefox\\firefox.exe", url)
		err := command.Run()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}
