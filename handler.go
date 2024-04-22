package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"magazine/api"
	"magazine/config"
	"magazine/db"
	"magazine/engine"
	"magazine/entity"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
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
		if homepage(ctx) {
			sendFileContent("./static/html/homepage.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "homepage" {
		if homepage(ctx) {
			sendFileContent("./static/html/homePage.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "login" {
		l := r.FormValue("username")
		p := r.FormValue("password")

		var user entity.User
		db.DB().Where("login = ? and password = ?", l, p).First(&user)

		if user.Uid != 0 {
			token := entity.Token{
				Uid:     user.Uid,
				Token:   api.CreateToken(user.Login, user.Password, 2*time.Minute),
				Expired: time.Now().Add(2 * time.Minute),
			}

			cookie := http.Cookie{
				Name:  "token",
				Value: token.Token,
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if pathArr[0] == "admin" {
		sendFileContent("./static/html/admin.html", ctx)
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

func homepage(ctx engine.Context) bool {
	cookie, err := ctx.Request.Cookie("token")
	if err == nil {
		// Проверка валидности токена
		if api.IsValidateToken(cookie.Value) {
			// Если токен валиден, перенаправляем пользователя на главную страниц
			return true
		}
	}
	return false
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
