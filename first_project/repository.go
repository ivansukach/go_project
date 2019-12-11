package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//CREATE
func (r repository) saveProduct(c echo.Context) error{
	p := Product{}
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
	r.InsertIntoDB(p)
	return c.String(http.StatusOK, 	"")
}
//DELETE and UPDATE
func (r repository) deleteProduct(c echo.Context) error {
	id := c.FormValue("id_delete")
	amountStr := c.FormValue("amount_delete")
	amount, _ := strconv.Atoi(amountStr)
	row, err := r.db.Query("SELECT amount FROM products WHERE id=?", id)
	amountRecord:=0
	defer row.Close()
	for row.Next() {
		err = row.Scan(&amountRecord)
	if amountRecord == amount {
	_, err = r.db.Query("DELETE * FROM products WHERE id=?", id)
	if err != nil {
		log.Println("Не получилось удалить данные")
		log.Println(err)
	}}else {
		if amountRecord > amount{
			_, err = r.db.Query("UPDATE products SET amount=? WHERE id=?", amountRecord-amount, id)
			if err != nil {
				log.Println("Не получилось удалить данные")
				log.Println(err)
			}
		}else{
			return c.String(http.StatusOK, "К сожалению Вы что-то перепутали")
		}
	}}
	log.Println(amountRecord)
	log.Println("Сработал Get-запрос на удаление")

	return c.String(http.StatusOK, "Удалено "+amountStr+" единиц товара с артикулом "+id)
}
//READ
func (r repository) getValuesFromDB(c echo.Context) error{
	log.Println("Сейчас будем получать данные из базы данных" )
	rows, err := r.db.Query("SELECT * FROM products")
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
func (r repository) InsertIntoDB(p Product){
	log.Println(p.Id)
	_, err := r.db.Query("INSERT INTO products VALUES (?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.Category, p.Date, p.Price, p.Amount)
	if err != nil {
		log.Println("В базу данных не смог вставить")
		log.Println(err)
	}
}
func (r repository) closeDB() error {
	r.db.Close()
	log.Println("Закрываем БД")
	return nil
}
type Repository interface {
	getValuesFromDB(c echo.Context) error
	saveProduct(c echo.Context) error
	deleteProduct(c echo.Context) error
	closeDB() error
}
type repository struct {
	db *sql.DB
}
