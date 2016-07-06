package gobotnet

import (
	"errors"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDir(dirPath string, fileMode os.FileMode) bool {
	err := os.MkdirAll(dirPath, fileMode)
	if err != nil {
		OutMessage(err.Error())
		return false
	}
	return true
}

func CreateFile(pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	defer file.Close()
	return nil
}

func WriteDataToFile(filePath string, data []byte) error {
	err := ioutil.WriteFile(filePath, data, 0644)
	OutMessage(err.Error())
	return err
}

func ReadFile(nameFile string) (bytesFile []byte, err error) {
	readBytes, err := ioutil.ReadFile(nameFile)
	OutMessage(err.Error())
	return readBytes, err
}

func SaveImageToFile(image *image.RGBA, nameFile string) error {
	f, err := os.Create("./" + nameFile)
	if err != nil {
		OutMessage(err.Error())
		return err
	}

	err = png.Encode(f, image)
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	f.Close()

	return nil
}

func CopyFileToDirectory(pathSourceFile string, pathDestFile string) error {
	sourceFile, err := os.Open(pathSourceFile)
	if CheckError(err) {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if CheckError(err) {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if CheckError(err) {
		return err
	}

	err = destFile.Sync()
	if CheckError(err) {
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if CheckError(err) {
		return err
	}

	destFileInfo, err := destFile.Stat()
	if CheckError(err) {
		return err
	}

	if sourceFileInfo.Size() == destFileInfo.Size() {
	} else {
		err = errors.New("Bad copy file")
		CheckError(err)
		return err
	}
	return nil
}

func DeleteFile(nameFile string) error {
	err := os.Remove(nameFile)
	CheckError(err)
	return err
}

func RemoveDirWithContet(dir string) error {
	d, err := os.Open(dir)
	if CheckError(err) {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if CheckError(err) {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if CheckError(err) {
			return err
		}
	}
	err = os.RemoveAll(dir)
	if CheckError(err) {
		return err
	}
	return nil
}

func SaveToken(pathFile, token string) bool {
	err := ioutil.WriteFile(pathFile, []byte(token), 0644)
	return CheckError(err)
}

func LoadToken(pathFile string) string {
	readBytes, err := ioutil.ReadFile(tokenFile)
	if CheckError(err) {
		return ""
	}
	return string(readBytes)
}
