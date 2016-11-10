package wrapper

import (
	"syscall"
	"unsafe"
)

var (
	winhttpdll, _                            = syscall.LoadLibrary("Winhttp.dll")
	kernel32dll, _                           = syscall.LoadLibrary("Kernel32.dll")
	winHttpGetIEProxyConfigForCurrentUser, _ = syscall.GetProcAddress(winhttpdll, "WinHttpGetIEProxyConfigForCurrentUser")

	winHttpAddRequestHeaders, _  = syscall.GetProcAddress(winhttpdll, "WinHttpAddRequestHeaders")
	winHttpOpen, _               = syscall.GetProcAddress(winhttpdll, "WinHttpOpen")
	winHttpConnect, _            = syscall.GetProcAddress(winhttpdll, "WinHttpConnect")
	winHttpOpenRequest, _        = syscall.GetProcAddress(winhttpdll, "WinHttpOpenRequest")
	winHttpSendRequest, _        = syscall.GetProcAddress(winhttpdll, "WinHttpSendRequest")
	winHttpReceiveResponse, _    = syscall.GetProcAddress(winhttpdll, "WinHttpReceiveResponse")
	winHttpCloseHandle, _        = syscall.GetProcAddress(winhttpdll, "WinHttpCloseHandle")
	winHttpReadData, _           = syscall.GetProcAddress(winhttpdll, "WinHttpReadData")
	winHttpQueryDataAvailable, _ = syscall.GetProcAddress(winhttpdll, "WinHttpQueryDataAvailable")
	winHttpSetOption, _          = syscall.GetProcAddress(winhttpdll, "WinHttpSetOption")
	winHttpSetTimeOut, _         = syscall.GetProcAddress(winhttpdll, "WinHttpSetTimeouts")
	getLastError, _              = syscall.GetProcAddress(kernel32dll, "GetLastError")
)

const (
	GET_REQUEST                     = "GET"
	POST_REQUEST                    = "POST"
	INTERNET_DEFAULT_HTTP_PORT      = 80
	INTERNET_DEFAULT_HTTPS_PORT     = 443
	WINHTTP_OPTION_PROXY            = 38
	WINHTTP_ACCESS_TYPE_NAMED_PROXY = 3
	WINHTTP_ADDREQ_FLAG_ADD         = 0x20000000
)

type winHttpIEProxyConfig struct {
	fAutoDetect       bool
	lpszAutoConfigUrl *uint32
	lpszProxy         *uint32
	lpszProxyBypass   *uint32
}

type winHttpProxyInfo struct {
	dwAccessType    uint32
	lpszProxy       *uint32
	lpszProxyBypass *uint32
}

func WinHttpGetIEProxyConfigForCurrentUser(lpIeProxy *winHttpIEProxyConfig) (result int) {
	ret, _, _ := syscall.Syscall(uintptr(winHttpGetIEProxyConfigForCurrentUser),
		1,
		uintptr(unsafe.Pointer(lpIeProxy)),
		0,
		0)
	return int(ret)
}

func WinHttpOpen(useragent string) (resp uintptr) {
	hSession, _, _ := syscall.Syscall6(uintptr(winHttpOpen),
		5,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(useragent))),
		0, 0, 0, 0, 0)
	return hSession
}

func WinHttpConnect(hSession uintptr, url string, port int) (resp uintptr) {
	hConnect, _, _ := syscall.Syscall6(uintptr(winHttpConnect),
		4,
		hSession,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(url))),
		uintptr(port),
		0, 0, 0)
	return hConnect
}

func WinHttpOpenRequest(hConnect uintptr, method string, uri string) (resp uintptr) {
	hRequest, _, _ := syscall.Syscall9(uintptr(winHttpOpenRequest),
		7,
		hConnect,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(method))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(uri))),
		0, 0, 0, 0, 0, 0)
	return hRequest
}

func WinHttpSetTimeOut(hSession uintptr, dwResolveTimeout, dwConnectTimeout, dwSendTimeout, dwReceiveTimeout int) (response uintptr) {
	result, _, _ := syscall.Syscall9(uintptr(winHttpSetTimeOut), 5,
		hSession,
		uintptr(dwResolveTimeout),
		uintptr(dwConnectTimeout),
		uintptr(dwSendTimeout),
		uintptr(dwReceiveTimeout), 0, 0, 0, 0)
	return result
}

func WinHttpSendRequest(hRequest uintptr, optional []byte, optional_len int) (resp uintptr) {
	var bResults uintptr
	if optional_len == 0 {
		bResults, _, _ = syscall.Syscall9(uintptr(winHttpSendRequest), 7,
			hRequest, 0, 0, 0, 0, 0, 0, 0, 0)
	} else {
		bResults, _, _ = syscall.Syscall9(uintptr(winHttpSendRequest), 7,
			hRequest, 0, 0,
			uintptr(unsafe.Pointer(&optional[0])),
			uintptr(optional_len),
			uintptr(optional_len),
			0, 0, 0)
	}
	return bResults
}

func WinHttpReceiveResponse(hRequest uintptr) (resp uintptr) {
	bResults, _, _ := syscall.Syscall(uintptr(winHttpReceiveResponse), 2, hRequest, 0, 0)
	return bResults
}

func WinHttpQueryDataAvailable(hRequest uintptr, dwSize *uint32) {
	syscall.Syscall(uintptr(winHttpQueryDataAvailable), 2, hRequest, uintptr(unsafe.Pointer(dwSize)), 0)
}

func WinHttpReadData(hRequest uintptr, buf []byte, dwSize uint32, dwDownloaded *uint32) {
	syscall.Syscall6(uintptr(winHttpReadData),
		4,
		hRequest,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(dwSize),
		uintptr(unsafe.Pointer(&*dwDownloaded)), 0, 0)
}

func WinHttpAddRequestHeaders(hRequest uintptr, headers string, flag int) {
	syscall.Syscall6(uintptr(winHttpAddRequestHeaders),
		4,
		hRequest,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(headers))),
		uintptr(len(headers)),
		uintptr(flag), 0, 0)
}

func WinHttpSetOption(hSession uintptr, option uintptr, buffer uintptr, buffersize uintptr) {
	syscall.Syscall6(uintptr(winHttpSetOption),
		4,
		hSession,
		option,
		buffer,
		buffersize, 0, 0)
}

func WinHttpCloseHandle(hHandle uintptr) {
	syscall.Syscall(uintptr(winHttpCloseHandle), 1, hHandle, 0, 0)
}
