package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)
type ViewData struct{
	Title string
	Message string
}
func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
var database *sql.DB
//type repository struct {
//	db *sql.DB
//}
//
//type Repository interface {
//	Create() error
//	Get() ()
//}
//func OpenMysqlRepository(cfg *config.Config) (Repository, error) {
//
//	//db, err := sql.Open("mysql", cfg.ConnectionString)
//	db, err := sql.Open("mysql", "root:6854321a@/skis&snowboard_shop")
//	if err != nil {
//		return nil, err
//	}
//
//	return &repository{db: db}, nil
//
//}JWTConfig struct {
////	Skipper Skipper
////	SigningKey interface{}
////	SigningMethod string
////	ContextKey string
////	Claims jwt.Claims
////	TokenLookup string
////	AuthScheme string
////}
////DefaultJWTConfig = JWTConfig{
////	Skipper: DefaultSkipper,
////	SigningMethod: AlgorithmHS256,
////	ContextKey: "user",
////	TokenLookup: "header:" + echo.HeaderAuthorization,
////	AuthScheme: "Bearer",
////	Claims: jwt.MapClaims{},
////}
//
type Product struct{
	Id string
	Name string
	Category string
	Date string
	Price float64
	Amount int
}

func IndexHandler(c echo.Context) error{
	log.Println("Сейчас будем получать данные из базы данных" )
	//rows, err := database.Query("SELECT * FROM skis&snowboard_shop.products")
	rows, err := database.Query("SELECT * FROM products")
	if err != nil {
		log.Println("Ошибка: не могу получить данные")
		log.Println(err)
	}
	defer rows.Close()
	products := []Product{}
	log.Println("Работаем дальше")
	for rows.Next(){
		p := Product{}
		err := rows.Scan(&p.Id, &p.Name, &p.Category, &p.Date, &p.Price, &p.Amount)
		if err != nil{
			fmt.Println(err)
			log.Println("Ошибочка в сканировании")
			continue
		}
		products = append(products, p)
	}
	tmpl, _ := template.ParseFiles("templates/db_example.html")
	tmpl.Execute(c.Response().Writer, products)
	return nil
}
func StatusHandler (c echo.Context) error{
	c.Response().Writer.Write([]byte("API IS RUNNING"))
	return nil
}


var exampleSetOfProducts = []Product{
	Product{Id: "hover-shooters", Name: "Hover Shooters", Category: "hover-shooters",
	Date : "2019-11-11", Price : 200.1, Amount : 2},
	Product{Id: "ocean-explorer", Name: "Ocean Explorer", Category: "hover-shooters",
	Date : "2019-11-11", Price : 200.2, Amount : 2},
	Product{Id: "dinosaur-park", Name: "Dinosaur Park", Category: "hover-shooters",
	Date : "2019-11-11", Price : 200.3, Amount : 2},
	Product{Id: "cars-vr", Name: "Cars VR", Category: "hover-shooters",
	Date: "2019-11-11", Price : 200.4, Amount : 2},
	Product{Id: "robin-hood", Name: "Robin Hood", Category: "hover-shooters" ,
	Date : "2019-11-11", Price : 200.5, Amount : 2},
	Product{Id: "real-world-vr", Name: "Real World VR", Category: "hover-shooters" ,
	Date : "2019-11-11", Price : 200.6, Amount : 2},
}

func ProductsHandler (c echo.Context) error{
	payload, _ := json.Marshal(exampleSetOfProducts)

	c.Response().Writer.Header().Set("Content-Type", "application/json")
	c.Response().Writer.Write([]byte(payload))
	return nil
}
func AddFeedbackHandler (c echo.Context) error{
	var product Product
	r:=c.Request()
	vars := mux.Vars(r)
	id := vars["id"]

	for _, p := range exampleSetOfProducts {
		if p.Id == id {
			product = p
			}
	}

	c.Response().Writer.Header().Set("Content-Type", "application/json")
	if product.Id != "" {
		payload, _ := json.Marshal(product)
		c.Response().Writer.Write([]byte(payload))
	} else {
		c.Response().Writer.Write([]byte("Product Not Found"))
	}
	return nil
}
func InsertIntoDB(/*w http.ResponseWriter, r *http.Request, */p Product){
	log.Println(p.Id)
	//_, err := database.Query("INSERT INTO skis_and_snowboard_shop.products VALUES (&p.id, &p.name, &p.category, &p.date, &p.price, &p.amount)")
	//_, err := database.Query("INSERT INTO products VALUES ($1, $2, $3, $4, $5, $6)", "1332", "uladzimir", "ivanukovich", "2019-11-11", 222, 11)
	_, err := database.Query("INSERT INTO products VALUES (?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.Category, p.Date, p.Price, p.Amount)
	//_, err := database.Query("INSERT INTO skis_and_snowboard_shop.products VALUES (?, ?, ?, ?, ?, ?)", "1332", "uladzimir", "ivanukovich", "2019-11-11", 222.2, 11)
	if err != nil {
		log.Println("В базу данных не смог вставить")
		log.Println(err)
	}

}
func saveProduct(c echo.Context) error{
	p := Product{}
	log.Println("В графе amount: " + c.FormValue("amount") )
	p.Id = c.FormValue("id")
	p.Name = c.FormValue("name")
	p.Category = c.FormValue("category")
	p.Date = c.FormValue("date")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	log.Println("В графе amount: " + c.FormValue("amount") )
	amount, err := strconv.Atoi(c.FormValue("amount") )
	if err != nil {
		log.Println("Ошибка")
		log.Println(err)
	}
	p.Price = price
	p.Amount = amount

	log.Println("id: "+p.Id)
	log.Println("name: "+p.Name)
	log.Println("category: "+p.Category)
	log.Println("date: "+p.Date)
	log.Print("price: ")
	log.Println(p.Price)
	log.Print("amount: ")
	log.Println(p.Amount)
	InsertIntoDB(p)
	return c.String(http.StatusOK, 	"")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "ivan" || password != "1111" {
		return echo.ErrUnauthorized
	}
	token := jwt.New(jwt.SigningMethodHS256)

	// Устанавливаем набор параметров для токена
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "Ivanko Sukach"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Подписываем токен нашим секретным ключем
	mySigningKey := []byte("secret")
	tokenString, _ := token.SignedString(mySigningKey)

			data := ViewData{
				Title: tokenString,
			}
	tmpl, _ := template.ParseFiles("templates/sign-in.html")
	tmpl.Execute(c.Response().Writer, data)


	return nil
}
func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

//func restricted(c echo.Context) error {
//	user := c.Get("user").(*jwt.Token)
//	claims := user.Claims.(jwt.MapClaims)
//	name := claims["name"].(string)
//	return c.String(http.StatusOK, "Welcome "+name+"!")
//}
func restricted(c echo.Context) error {
	user:=c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	//r := mux.NewRouter()
	//fs :=http.FileServer(http.Dir("templates"))
	//http.Handle("/", fs)
	//log.Println("Listening ...")
	//conn, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/story")
	db, err := sql.Open("mysql", "root:6854321a@/skis_and_snowboard_shop")
	if err != nil {
		log.Println(err)
	}
	database = db
	//http.HandleFunc("/hello", IndexHandler)
	//http.ListenAndServe(":3000", nil)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//	SigningKey: []byte("secret"),
	//	TokenLookup: "query:token",
	//}))
	e.GET("/", accessible)
	e.Static("/", "templates")
	//e.POST("/post", saveProduct)
	e.POST("/sign-in", login)
	e.GET("/db_example", IndexHandler)
	e.GET("/status", StatusHandler)
	e.POST("/restricted", restricted)
	//e.GET("/products", ProductsHandler)
	//e.GET("/products/{id}/feedback", AddFeedbackHandler)
	//e.GET("/get-token", login)
	//e.GET("/get", func(c echo.Context) error {
	//		data := ViewData{
	//			Title: c.FormValue("id_delete"),
	//			Message: c.FormValue("amount_delete"),
	//		}
	//	wr := c.Response().Writer
	//	log.Println("Сработал Get-запрос")
	//	log.Println(data.Title)
	//	log.Println(data.Message)
	//	tmpl, _ :=
	//	template.ParseFiles("templates/login.html")
	//	tmpl.Execute(wr, data)
	//	return nil
	//})
	// Restricted group
	r := e.Group("/restricted")
	//r.Use(middleware.JWT([]byte("secret")))
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		TokenLookup: "query:token",
	}))
	r.GET("", restricted)

	defer db.Close()
	e.Logger.Fatal(e.Start(":1323"))
}