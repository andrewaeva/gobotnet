package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	gobotnetPath := os.Args[1]
	fullpath := os.Args[2]
	filename := filepath.Base(fullpath)

	input, err := os.Open(fullpath)

	if err != nil {
		fmt.Println(err)
		return
	}
	fi, err := input.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	b1 := make([]byte, fi.Size())
	_, err = input.Read(b1)

	if err != nil {
		fmt.Println(err)
		return
	}

	encoded_file := base64.StdEncoding.EncodeToString(b1)
	all_data := "package addfile\n\nimport (\n\"encoding/base64\"\n)\n\nconst (\n str_file = \"" + encoded_file + "\"\n filename=\"" + filename + "\"\n)\n\nfunc GetAdditionFile()(file []byte) {\n decoded,_ := base64.StdEncoding.DecodeString(str_file) \n return decoded\n}\n\nfunc GetAdditionFileName()(name string) {\nreturn filename\n}\n"
	output, err := os.Create(gobotnetPath + "src/addfile/addfile.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = output.Write([]byte(all_data))
	if err != nil {
		fmt.Println(err)
		return
	}

}
