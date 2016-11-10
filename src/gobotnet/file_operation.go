package gobotnet

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		CheckError(err)
		return false
	} else {
		return true
	}
}

func CreateDir(dirPath string, fileMode os.FileMode) bool {
	err := os.MkdirAll(dirPath, fileMode)
	if CheckError(err) {
		return false
	}
	return true
}

func CreateFile(pathFile string) error {
	file, err := os.Create(pathFile)
	if CheckError(err) {
		return err
	}
	defer file.Close()
	return nil
}

func CreateFileAndWriteData(fileName string, writeData []byte) error {
	fileHandle, err := os.Create(fileName)

	CheckError(err)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	writer.Write(writeData)
	writer.Flush()
	return nil
}

func ReadFile(nameFile string) (bytesFile []byte, err error) {
	readBytes, err := ioutil.ReadFile(nameFile)
	CheckError(err)
	return readBytes, err
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

func RemoveDirWithContent(dir string) error {
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
