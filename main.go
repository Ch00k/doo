package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

type Entry struct {
	gorm.Model
	CompletedAt time.Time
	Text        string
	Comments    []Comment
}

type Comment struct {
	gorm.Model
	EntryID uint
	Text    string
}

func initDB(host, port, username, password, dbname string) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, username, password, dbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
}

func getEnvVar(name, default_ string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	} else {
		return default_
	}
}

func setupDB() *gorm.DB {
	dbHost := getEnvVar("DOO_DB_HOST", "localhost")
	dbPort := getEnvVar("DOO_DB_PORT", "5432")
	dbUser := getEnvVar("DOO_DB_USER", "doo")
	dbPass := getEnvVar("DOO_DB_PASSWORD", "doo")
	dbName := getEnvVar("DOO_DB_NAME", "doo")

	db, err := initDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&Entry{},
		&Comment{},
	)

	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/entries", func(c *gin.Context) {
		var entries []Entry
		db.Unscoped().Preload(clause.Associations).Find(&entries)
		c.JSON(http.StatusOK, entries)
	})

	r.GET("/entries/:id", func(c *gin.Context) {
		var entry Entry
		res := db.Preload(clause.Associations).First(&entry, c.Param("id"))
		if res.RowsAffected < 1 {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, entry)
		}
	})

	r.POST("/entries", func(c *gin.Context) {
		var entry Entry
		if err := c.ShouldBindJSON(&entry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res := db.Create(&entry)
		if res.Error == nil {
			c.Status(http.StatusCreated)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	})

	r.PUT("/entries/:id", func(c *gin.Context) {
		var entry Entry
		res := db.First(&entry, c.Param("id"))
		if res.RowsAffected < 1 {
			c.Status(http.StatusNotFound)
		} else {
			var updatedEntry Entry
			if err := c.ShouldBindJSON(&updatedEntry); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			db.Model(&entry).Updates(updatedEntry)
		}
	})

	r.DELETE("/entries/:id", func(c *gin.Context) {
		var entry Entry
		res := db.First(&entry, c.Param("id"))
		if res.RowsAffected < 1 {
			c.Status(http.StatusNotFound)
		} else {
			db.Delete(&entry)
		}
	})

	r.POST("/entries/:id/comment", func(c *gin.Context) {
		var comment Comment
		var entry Entry
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res := db.First(&entry, c.Param("id"))
		if res.RowsAffected < 1 {
			c.Status(http.StatusNotFound)
		} else {
			err := db.Model(&entry).Association("Comments").Append(&comment)
			if err == nil {
				c.Status(http.StatusCreated)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
	})

	r.DELETE("/entries/:id/comments/:cid", func(c *gin.Context) {
		var entry Entry
		res := db.First(&entry, c.Param("eid"))
		if res.RowsAffected < 1 {
			c.Status(http.StatusNotFound)
		} else {
			var comment Comment
			res := db.First(&comment, c.Param("cid"))
			if res.RowsAffected < 1 {
				c.Status(http.StatusNotFound)
			} else {
				db.Delete(&comment)
			}
		}
	})

	return r
}

func main() {
	db := setupDB()
	r := setupRouter(db)

	httpHost := getEnvVar("DOO_HTTP_HOST", "localhost")
	httpPort := getEnvVar("DOO_HTTP_PORT", "8080")
	r.Run(fmt.Sprintf("%s:%s", httpHost, httpPort))
}
