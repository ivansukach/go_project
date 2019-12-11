package main

import (
	"fmt"
	"github.com/labstack/echo"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type User struct{
	Password string
	Name string
	Surname string
	Age int
	IsAdministrator bool
	IsModerator bool
	Discount float64
}
//var Repo = make(map[string]User)
type Repo map[string]User
type Repository2 interface {
	getValuesFromDB2(c echo.Context) error
	saveUser(c echo.Context) error
	getPassword(c echo.Context) error
	closeDB() error
}

//CREATE
func (r repository) saveUser(c echo.Context) error{
	p := Repo{}
	u := User{}
	u.Password = c.FormValue("password")
	u.Name = c.FormValue("name")
	u.Surname = c.FormValue("surname")
	u.Age, _ = strconv.Atoi(c.FormValue("age"))
	u.IsAdministrator = false;
	u.IsModerator = false;
	u.Discount = 0;
	p[c.FormValue("email")]=u
	r.InsertIntoDB2(c.FormValue("email"), p)
	return c.String(http.StatusOK, 	"")
}
//DELETE and UPDATE
func (r repository) getPassword(c echo.Context) error {
	username := c.FormValue("username")
	row, err := r.db.Query("SELECT password FROM users WHERE login=?", username)
	password :=""
	defer row.Close()
	for row.Next() {
		err = row.Scan(&password)
		if err != nil{
			fmt.Println(err)
			log.Println("Ошибочка при получении пароля")
			continue
		}
	}
	if err != nil{
		fmt.Println(err)
		log.Println("Ошибочка при получении пароля")
	}
	log.Println(password)
	log.Println("Сработал Get-запрос на получение пароля")

	return c.String(http.StatusOK, password)
}
//READ
func (r repository) getValuesFromDB2(c echo.Context) error{
	log.Println("Сейчас будем получать данные из базы данных Пользователи" )
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		log.Println("Ошибка: не могу получить данные")
		log.Println(err)
	}
	defer rows.Close()
	users := []Repo{}
	log.Println("Работаем дальше")
	for rows.Next(){
		key := ""
		p := Repo{}
		u := User{}
		err := rows.Scan(&key, &u.Name, &u.Surname, &u.Age, &u.IsAdministrator, &u.IsModerator, &u.Discount)
		if err != nil{
			fmt.Println(err)
			log.Println("Ошибочка в сканировании БД")
			continue
		}
		p[key]=u

		users = append(users, p)
	}
	tmpl, _ := template.ParseFiles("templates/db_example.html")
	tmpl.Execute(c.Response().Writer, users)
	return nil
}
func (r repository) InsertIntoDB2(key string, p Repo){
	_, err := r.db.Query("INSERT INTO users VALUES (?, ?, ?, ?, ?, ?, ?, ?)", key, p[key].Name, p[key].Surname, p[key].Age, p[key].IsAdministrator, p[key].IsModerator, p[key].Discount)
	if err != nil {
		log.Println("В базу данных не смог вставить")
		log.Println(err)
	}
}
