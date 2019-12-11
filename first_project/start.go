package main

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"log"
	"net/http"
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


func OpenMysqlRepository() (Repository, error){
	db, err := sql.Open("mysql", "root:6854321a@/skis_and_snowboard_shop")
	if err != nil {
		log.Println(err)
	}
	return &repository{db: db}, nil
}
type Product struct{
	Id string
	Name string
	Category string
	Date string
	Price float64
	Amount int
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

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "ivan" || password != "1111" {
		return echo.ErrUnauthorized
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "Ivanko Sukach"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

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


func restricted(c echo.Context) error {
	user:=c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	rps, _ := OpenMysqlRepository()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", accessible)
	e.Static("/", "templates")
	e.POST("/post", rps.saveProduct)
	e.POST("/sign-in", login)
	e.GET("/db_example", rps.getValuesFromDB)
	e.GET("/status", StatusHandler)
	e.POST("/restricted", restricted)
	//e.GET("/products", ProductsHandler)
	//e.GET("/products/{id}/feedback", AddFeedbackHandler)
	//e.GET("/get-token", login)
	e.GET("/get", rps.deleteProduct)
	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		TokenLookup: "query:token",
	}))
	r.GET("", restricted)
	defer rps.closeDB()
	e.Logger.Fatal(e.Start(":1323"))
}