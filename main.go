package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)


type Shortener struct {
	Id 			int
	LongURL 	string 	`db:"longURL" form:"longURL"`
	ShortURL  	string	`db:"shortURL" form:"shortURL"`
}


var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/urlshortener") //koneksi DB

func generateSlug() string {
	var chars = []rune("1qaz2wsx3edc4rfv5tgb6yhn7ujm8ik9ol0p") //menentukan karakter apa saja yg akan dirandom
	s := make([]rune, 5)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))] // random rune
	}
	return string(s)
}


func add(c *gin.Context) {

	ShortURL := generateSlug() //memanggil func generateSlug
	LongURL := c.PostForm("longURL") // mengambil variable longURL dari Form


	stmt, err := db.Prepare("insert into shortener (shortURL, longURL) values(?,?);") // syntax insert data mysql
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(ShortURL, LongURL)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()


	c.HTML(http.StatusOK, "result.html", gin.H{ 
		"message": fmt.Sprintf("201 Created"),
		"LongURL": fmt.Sprintf("%s",LongURL),
		"ShortURL": fmt.Sprintf("http://localhost:8000/go/%s",ShortURL),
		"URL": ShortURL,
	})
	
}

func routeIndex(c *gin.Context)  {
    var shortener Shortener

    shortURL := c.Param("shortURL")
    row := db.QueryRow("select id, shortURL, longURL from shortener where shortURL = ?;", shortURL)
    err = row.Scan(&shortener.Id, &shortener.ShortURL, &shortener.LongURL)

    if err != nil {
		c.JSON(http.StatusOK, nil)
	} else {
		c.Redirect(301, shortener.LongURL)
	}
}

func displayHTML(c *gin.Context)  {
	c.HTML(http.StatusOK, "home.html", nil)
}

func displayHTML2(c *gin.Context)  {
	c.HTML(http.StatusOK, "coba.html", nil)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.POST("/create", add)
	router.GET("/go/:shortURL", routeIndex)
	router.GET("/", displayHTML)
	router.GET("/2", displayHTML2)
	router.Run(":8000")	
}