package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/entries", func(c *gin.Context) {
		var entries []Entry
		db.Preload(clause.Associations).Find(&entries)
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
			// TODO: This returns Comments as nil instead of []
			c.JSON(http.StatusCreated, entry)
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
			c.JSON(http.StatusOK, entry)
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
				c.JSON(http.StatusCreated, comment)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
	})

	// TODO: Update comments?

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
