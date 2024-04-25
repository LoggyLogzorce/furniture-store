package main

import (
	"encoding/json"
	"fmt"
	"furniture_store/api"
	"furniture_store/config"
	"furniture_store/db"
	"furniture_store/engine"
	"furniture_store/entity"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Info struct {
	Name   string
	Access string
	Link   reflect.Value
}

var types map[string]bool
var hdl *api.Handler
var apiMap map[string]map[string]Info
var accessExceptions []string

func init() {
	cfg := config.Get()
	apiMap = make(map[string]map[string]Info)
	apiMap["POST"] = make(map[string]Info)
	apiMap["PUT"] = make(map[string]Info)
	apiMap["DELETE"] = make(map[string]Info)
	apiMap["GET"] = make(map[string]Info)
	maps := cfg.Api

	types = make(map[string]bool)
	types[".png"] = true
	types[".html"] = true
	types[".js"] = true
	types[".svg"] = true
	types[".css"] = true

	hdl = &api.Handler{}
	services := reflect.ValueOf(hdl)
	_struct := reflect.TypeOf(hdl)

	for methodNum := 0; methodNum < _struct.NumMethod(); methodNum++ {
		method := _struct.Method(methodNum)
		val, ok := maps[method.Name]
		if !ok {
			continue
		}
		if _, ok := apiMap[val.Method]; !ok {

		}
		apiMap[val.Method][val.Url] = Info{
			Name:   method.Name,
			Access: "",
			Link:   services.Method(methodNum),
		}
	}

	accessExceptions = cfg.List
}

func mainHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	if pathArr[0] == "" {
		validToken, role := homepage(ctx)
		if validToken && role {
			sendFileContent("./static/html/admin.html", ctx)
			return
		}
		if validToken && !role {
			sendFileContent("./static/html/homepage.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "homepage" {
		validToken, _ := homepage(ctx)
		if validToken {
			sendFileContent("./static/html/homePage.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "register" {
		sendFileContent("./static/html/register.html", ctx)
		return
	}

	if pathArr[0] == "login" {
		var user entity.User
		user.Login = r.FormValue("username")
		user.Password = r.FormValue("password")

		cookie, err, admin := api.UserRead(user)
		fmt.Println(admin)

		if err == true {
			http.SetCookie(ctx.Response, &cookie)
			switch admin {
			case true:
				http.Redirect(ctx.Response, ctx.Request, "/admin", http.StatusOK)

			case false:
				http.Redirect(ctx.Response, ctx.Request, "/homepage", http.StatusOK)
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if pathArr[0] == "logout" {
		var user entity.User
		user.Name = r.FormValue("username")
		user.Login = r.FormValue("login")
		user.Password = r.FormValue("password")
		user.Addm = false

		err := api.CreateUser(user)
		if err == true {
			http.Redirect(ctx.Response, ctx.Request, "/", http.StatusOK)
			return
		}
		http.Error(w, "Unregistered", http.StatusConflict)
		return
	}

	if pathArr[0] == "admin" {
		validToken, role := homepage(ctx)
		if validToken && role {
			sendFileContent("./static/html/admin.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "data" {
		var users []entity.User
		db.DB().Find(&users)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(users)
		if err != nil {
			return
		}
	}

	_, ok := apiMap[r.Method]
	if !ok {
		http.Error(w, "No such method", http.StatusNotFound)
		return
	}

	if staticUrl, ok := static(path); ok {
		sendFileContent("./static/"+staticUrl, ctx)
		return
	}
}

func sendFileContent(filename string, ctx engine.Context) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		http.Error(ctx.Response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Определите Content-Type на основе расширения файла с помощью пакета mime
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream" // По умолчанию, если тип не найден
	}
	ctx.Response.Header().Set("Content-Type", contentType)

	_, err = ctx.Response.Write(file)
	if err != nil {
		return
	}
}

func homepage(ctx engine.Context) (bool, bool) {
	cookie, err := ctx.Request.Cookie("token")
	if err == nil {
		// Проверка валидности токена
		if api.IsValidateToken(cookie.Value) {
			role := api.GetRoleFromToken(cookie.Value)
			return true, role
		}
	}
	return false, false
}

func static(path string) (string, bool) {
	splitPath := strings.Split(path, "/")
	fileName := splitPath[len(splitPath)-1]
	splitName := strings.Split(fileName, ".")
	fileExt := "." + splitName[len(splitName)-1]
	if types[fileExt] {
		return path, true
	}
	return "", false
}

func sendFile(path string, ctx engine.Context) {
	http.ServeFile(ctx.Response, ctx.Request, path)
}

func updateHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	// Декодируем JSON данные из тела запроса
	var updatedData map[string]string
	if err := json.NewDecoder(ctx.Request.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	uintVal, err := strconv.ParseUint(updatedData["uid"], 10, 32)
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	Uid := uint32(uintVal)
	Addm, err := strconv.ParseBool(updatedData["Additional Permission"])
	if err != nil {
		fmt.Println("Ошибка конвертации:", err)
		return
	}
	user := entity.User{
		Uid:      Uid,
		Name:     updatedData["Name"],
		Login:    updatedData["Login"],
		Password: updatedData["Password"],
		Addm:     Addm,
	}

	db.DB().Save(&user)
}
