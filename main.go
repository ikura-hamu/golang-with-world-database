package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int    `json:"id,omitempty" db:"ID"`
	Name        string `json:"name,omitempty" db:"Name"`
	CountryCode string `json:"countryCode,omitempty" db:"CountryCode"`
	District    string `json:"district,omitempty" db:"District"`
	Population  int    `json:"population,omitempty" db:"Population"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	e := echo.New()

	e.GET("/city/:cityName", getCityInformHandler)

	e.POST("/addCity", addCity)

	e.Start(":11000")
}

func getCityInformHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	var city City
	if err := db.Get(&city, "select * from city where Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
		log.Printf("No such city name = %s", cityName)
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	return c.JSON(http.StatusOK, city)
}

func addCity(c echo.Context) error {
	newCity := &City{}
	err := c.Bind(newCity)

	if err != nil { //when an error happened
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%+v", newCity))
	}
	_, er := db.Exec(fmt.Sprintf("INSERT INTO city (Name, CountryCode, District, Population) VALUES ('%v', '%v', '%v', '%d')", newCity.Name, newCity.CountryCode, newCity.District, newCity.Population))
	if er != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("afterSQL--%+v", er))
	}
	return c.JSON(http.StatusOK, "added successfully")
}
