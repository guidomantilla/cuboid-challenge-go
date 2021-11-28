package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListBags(context *gin.Context) {
	var bags []models.Bag
	if r := db.CONN.Preload("Cuboids").Find(&bags); r.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	context.JSON(http.StatusOK, bags)
}

func GetBag(context *gin.Context) {
	bagID := context.Param("bagID")

	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	context.JSON(http.StatusOK, &bag)
}

func CreateBag(context *gin.Context) {
	var bagInput struct {
		Title  string
		Volume uint
	}

	if err := context.BindJSON(&bagInput); err != nil {
		return
	}

	bag := models.Bag{
		Title:   bagInput.Title,
		Volume:  bagInput.Volume,
		Cuboids: []models.Cuboid{},
	}
	if r := db.CONN.Create(&bag); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	context.JSON(http.StatusCreated, &bag)
}

func DeleteBag(context *gin.Context) {
	bagID := context.Param("bagID")

	var bag models.Bag
	if r := db.CONN.First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&bag); r.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "OK"})
}
