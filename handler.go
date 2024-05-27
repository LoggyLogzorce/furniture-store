package main

import (
	"encoding/json"
	"fmt"
	"furniture_store/api"
	"furniture_store/db"
	"furniture_store/engine"
	"furniture_store/entity"
	"io/ioutil"
	"log"
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
var apiMap map[string]map[string]Info

func init() {
	apiMap = make(map[string]map[string]Info)
	apiMap["POST"] = make(map[string]Info)
	apiMap["PUT"] = make(map[string]Info)
	apiMap["DELETE"] = make(map[string]Info)
	apiMap["GET"] = make(map[string]Info)

	types = make(map[string]bool)
	types[".png"] = true
	types[".html"] = true
	types[".js"] = true
	types[".svg"] = true
	types[".css"] = true
}

func userHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	if pathArr[0] == "" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken {
			switch role {
			case "admin":
				sendFileContent("./static/html/adminUsers.html", ctx)
				return
			case "user":
				sendFileContent("./static/html/homePage.html", ctx)
				return
			}
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "homepage" {
		validToken, _ := validateTokenAndRole(ctx)
		if validToken {
			sendFileContent("./static/html/homePage.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "order" {
		validToken, _ := validateTokenAndRole(ctx)
		if validToken {
			sendFileContent("./static/html/order.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[0] == "personal_account" {
		validToken, _ := validateTokenAndRole(ctx)
		if validToken {
			sendFileContent("./static/html/personal_account.html", ctx)
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

		cookie, role := api.UserRead(user)

		response := struct {
			Role string `json:"role"`
		}{
			Role: role,
		}

		ctx.Response.Header().Set("Content-Type", "application/json")
		if cookie.Value != "" {
			http.SetCookie(ctx.Response, &cookie)
			ctx.Response.WriteHeader(http.StatusOK)
			err := json.NewEncoder(ctx.Response).Encode(response)
			if err != nil {
				log.Println(err)
			}
		} else {
			http.Error(ctx.Response, "Unauthorized", http.StatusUnauthorized)
		}
	}

	if pathArr[0] == "logout" {
		var user entity.User
		user.Name = r.FormValue("username")
		user.Login = r.FormValue("login")
		user.Password = r.FormValue("password")
		user.Role = "user"

		created := api.CreateUser(user)
		if created == true {
			http.Redirect(ctx.Response, ctx.Request, "/", http.StatusOK)
			return
		}
		http.Error(w, "Unregistered", http.StatusConflict)
		return
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

func validateTokenAndRole(ctx engine.Context) (bool, string) {
	cookie, err := ctx.Request.Cookie("token")
	if err == nil {
		// Проверка валидности токена
		if api.IsValidateToken(cookie.Value) {
			role, e := api.GetRoleFromToken(cookie.Value)
			if e != nil {
				log.Println(e)
			}
			return true, role
		}
	}
	return false, ""
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

func updateHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	validToken, role := validateTokenAndRole(ctx)
	if validToken && role != "admin" {
		return
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	// Декодируем JSON данные из тела запроса
	var updatedData map[string]string
	if err := json.NewDecoder(ctx.Request.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if pathArr[2] == "user" {
		api.UpdateUser(updatedData)
		return
	}

	if pathArr[2] == "product" {
		api.UpdateProduct(updatedData)
		return
	}

	if pathArr[2] == "category" {
		api.UpdateCategory(updatedData)
		return
	}

	if pathArr[2] == "itemOrder" {
		api.UpdateItemsOrder(updatedData)
		return
	}
}

func deleteHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	var rowData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&rowData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if pathArr[2] == "item" {
		fmt.Println(rowData)
		cookie, err := ctx.Request.Cookie("token")
		if err != nil {
			log.Println(err)
		}

		token := cookie.Value
		api.DeleteItemsOrder(rowData, token)
		response := struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println(err)
		}
		return
	}

	validToken, role := validateTokenAndRole(ctx)
	if validToken && role != "admin" {
		return
	}

	if pathArr[2] == "review" {
		api.DeleteReview(rowData)
		return
	}

	if pathArr[2] == "user" {
		api.DeleteUser(rowData)
		return
	}

	if pathArr[2] == "product" {
		api.DeleteProduct(rowData)
		return
	}

	if pathArr[2] == "category" {
		api.DeleteCategory(rowData)
		return
	}
}

func addHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	var rowData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&rowData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if pathArr[2] == "cart" {
		cookie, err := ctx.Request.Cookie("token")
		if err != nil {
			log.Println(err)
		}

		token := cookie.Value
		api.AddCart(rowData, token)
		response := struct {
			Success bool `json:"success"`
		}{
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println(err)
		}
		return
	}

	validToken, role := validateTokenAndRole(ctx)
	if validToken && role != "admin" {
		return
	}

	if pathArr[2] == "product" {
		api.AddProduct(rowData)
		return
	}

	if pathArr[2] == "category" {
		api.AddCategory(rowData)
		return
	}
}

func getDataHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	if pathArr[1] == "order" {
		cookie, err := ctx.Request.Cookie("token")
		login, err := api.GetLoginFromToken(cookie.Value)
		var orderItem []entity.OrderItem
		db.DB().Table("product").
			Select("product.id as id, product.name as product_name, product.price as price, items_order.quantity as quantity").
			Joins("JOIN items_order ON product.id = items_order.product").
			Joins("JOIN \"order\" ON items_order.order_id = \"order\".id").
			Joins("JOIN \"user\" ON \"order\".user_id = \"user\".uid").
			Where("\"user\".login = ?", login).Scan(&orderItem)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(orderItem)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if pathArr[1] == "users" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			var users []entity.User
			db.DB().Find(&users)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(users)
			if err != nil {
				log.Println(err)
			}
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if pathArr[1] == "itemsOrders" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			var itemsOrderProducts []entity.ItemOrderProduct
			db.DB().Table("items_order").
				Select("items_order.id, items_order.order_id, product.id as product_id, product.name AS product_name, items_order.quantity, product.price").
				Joins("JOIN product ON items_order.product = product.id").
				Scan(&itemsOrderProducts)

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(itemsOrderProducts)
			if err != nil {
				log.Println(err)
			}
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if pathArr[1] == "products" {
		var products []entity.Product
		db.DB().Find(&products)
		categoryStr := r.URL.Query().Get("category")

		// Логируем значение параметра для отладки
		fmt.Printf("Запрошенная категория: %q\n", categoryStr)

		if categoryStr == "" {
			// Если категория не указана, возвращаем все товары
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
			return
		}

		// Преобразуем строку в uint8
		category, err := strconv.ParseUint(categoryStr, 10, 8)
		if err != nil {
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		fmt.Printf("Преобразованная категория: %d\n", category)
		category32 := uint32(category)

		var filteredProducts []entity.Product
		for _, product := range products {
			if product.Category == category32 {
				filteredProducts = append(filteredProducts, product)
			}
		}

		fmt.Printf("Отфильтрованные товары: %+v\n", filteredProducts)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredProducts)
		return
	}

	if pathArr[1] == "categories" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			var category []entity.Category
			db.DB().Find(&category)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(category)
			if err != nil {
				log.Println(err)
			}
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if pathArr[1] == "reviews" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			var review []entity.Review
			db.DB().Find(&review)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(review)
			if err != nil {
				log.Println(err)
			}
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
}

func adminHandle(w http.ResponseWriter, r *http.Request) {
	ctx := engine.Context{
		Response: w,
		Request:  r,
	}

	url := r.URL
	path := url.Path[1:]
	pathArr := strings.Split(path, "/")

	if pathArr[1] == "users" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			sendFileContent("./static/html/adminUsers.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[1] == "products" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			sendFileContent("./static/html/adminProducts.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[1] == "categories" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			sendFileContent("./static/html/adminCategory.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[1] == "itemsOrders" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			sendFileContent("./static/html/adminItemsOrders.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
	}

	if pathArr[1] == "reviews" {
		validToken, role := validateTokenAndRole(ctx)
		if validToken && role == "admin" {
			sendFileContent("./static/html/adminReviews.html", ctx)
			return
		}
		sendFileContent("./static/html/index.html", ctx)
		return
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
