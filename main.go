package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	//"os"
	//"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
		
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")

	//db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	
	err = db.Ping()
	if err != nil {
		panic("Gagal Menghubungkan ke Database")
	}
	defer db.Close()

	router := gin.Default()

	type About struct {

	//	gorm.Model
		
		Id    			int    `json: "id"`
		Name  		string `json: "name"`
		Image 			string `json: "image"`
		
	}

	// GET all persons
	router.GET("/get-about", func(c *gin.Context) {
		var (
			about  About
			abouts []About
		)
		rows, err := db.Query("select id, name, image from profile;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&about.Id, &about.Name, &about.Image)
			abouts = append(abouts, about)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"return":  abouts,
			"count": len(abouts),
		})

		// var abouts About
		// db.find(abouts);
		// c.JSON(200, gin.H{
		// 	"status": 200,
		// 	"data":   abouts,
		// })

	})

	// Menampilkan Detail Data Berdasarkan ID
	router.GET("/get-about/:id", func(c *gin.Context) {
		var (
			about  About
			result gin.H
		)
		ids := c.Param("id")
		row := db.QueryRow("select id, name, image from profile where id = ?;", ids)
		err = row.Scan(&about.Id, &about.Name, &about.Image)
		if err != nil {
			// If no results send null
			fmt.Print(err.Error())
			result = gin.H{
				"return": "Tidak ada data about yang ditemukan",
			}
		} else {
			result = gin.H{
				"return":  about,
				"count": 1,
			}
		}
		c.JSON(http.StatusOK, result)
	})


	// POST new person details
	router.POST("/get-about", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		image := c.PostForm("image")
		
		stmt, err := db.Prepare("insert into profile (id, name, image) values(?,?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id, name, image)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(image)
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprintf(" Berhasil menambahkan data %s ", datanya),
		})
	})

	// PUT - update a person details
	router.PUT("/get-about", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.PostForm("id")
		name := c.PostForm("name")
		image := c.PostForm("image")
		
		stmt, err := db.Prepare("update profile set name= ?, image = ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(name, image, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(image)
		
		defer stmt.Close()
		datanya := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"Pesannya": fmt.Sprintf("Berhasil Merubah Id %s Menjadi %s", id, datanya),
		})
	})

	// Delete resources
	router.DELETE("/get-about", func(c *gin.Context) {
		id := c.PostForm("id")
		stmt, err := db.Prepare("delete from profile where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"Pesannya": fmt.Sprintf("Berhasil Menghapus %s", id),
		})
	})

	router.Run(":80")
}
