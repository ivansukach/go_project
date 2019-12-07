package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
//}
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

func main() {
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
	e.Static("/", "templates")
	e.POST("/post", saveProduct)
	e.GET("/db_example", IndexHandler)
	e.GET("/get", func(c echo.Context) error {
			data := ViewData{
				Title: c.FormValue("id_delete"),
				Message: c.FormValue("amount_delete"),
			}
		wr := c.Response().Writer
		log.Println("Сработал Get-запрос")
		log.Println(data.Title)
		log.Println(data.Message)
		tmpl, _ :=
		template.ParseFiles("templates/login.html")
		tmpl.Execute(wr, data)
		return nil
	})

	defer db.Close()
	e.Logger.Fatal(e.Start(":1323"))
}