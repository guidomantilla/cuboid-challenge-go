package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListCuboids(context *gin.Context) {
	var cuboids []models.Cuboid
	if r := db.CONN.Find(&cuboids); r.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	context.JSON(http.StatusOK, cuboids)
}

func GetCuboid(context *gin.Context) {
	cuboidID := context.Param("cuboidID")

	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	context.JSON(http.StatusOK, &cuboid)
}

func CreateCuboid(context *gin.Context) {
	var cuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := context.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboid := models.Cuboid{
		Width:  cuboidInput.Width,
		Height: cuboidInput.Height,
		Depth:  cuboidInput.Depth,
		BagID:  cuboidInput.BagID,
	}

	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if bag.Disable {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})

		return
	}

	if cuboid.PayloadVolume() > bag.AvailableVolume() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})

		return
	}

	if r := db.CONN.Create(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	context.JSON(http.StatusCreated, &cuboid)
}

func UpdateCuboid(context *gin.Context) {
	cuboidID := context.Param("cuboidID")

	var cuboid models.Cuboid
	if r := db.CONN.Preload("Bag").First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	var cuboidInput struct {
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := context.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboid.Width = cuboidInput.Width
	cuboid.Height = cuboidInput.Height
	cuboid.Width = cuboidInput.Width
	cuboid.Depth = cuboidInput.Depth

	if cuboid.Bag.Disable {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})

		return
	}

	if cuboid.PayloadVolume() > cuboid.Bag.AvailableVolume() {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})

		return
	}

	if r := db.CONN.Updates(&cuboid); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	context.JSON(http.StatusOK, &cuboid)
}

func DeleteCuboid(context *gin.Context) {
	cuboidID := context.Param("cuboidID")

	var cuboid models.Cuboid
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&cuboid); r.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "OK"})
}
