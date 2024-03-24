package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Store godoc
// @Summary      Create an photo
// @Description  create and store an photo
// @Tags         Photo
// @Accept       json
// @Produce      json
// @Param        json  body  models.Photo true  "Photo"
// @Success      201  {object}  models.Photo
// @Security     Bearer
// @Router       /photos        [post]
func PhotoCreate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserId = userID

	err := db.Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"created_at": Photo.CreatedAt,
	})
}

// Fetch godoc
// @Summary      Fetch photos
// @Description  get photos
// @Tags         Photo
// @Accept       json
// @Produce      json
// @Success      200	{object}	[]models.Photo
// @Security     Bearer
// @Router       /photos        [get]
func PhotoGetAll(c *gin.Context) {
	db := database.GetDB()
	var Photos []models.Photo

	var data []interface{}

	err := db.Preload("User").Find(&Photos).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for i := range Photos {
		photo := make(map[string]interface{})
		user := make(map[string]interface{})

		user["email"] = Photos[i].User.Email
		user["username"] = Photos[i].User.Username

		photo["id"] = Photos[i].ID
		photo["title"] = Photos[i].Title
		photo["caption"] = Photos[i].Caption
		photo["photo_url"] = Photos[i].PhotoUrl
		photo["user_id"] = Photos[i].UserId
		photo["created_at"] = Photos[i].CreatedAt
		photo["updated_at"] = Photos[i].UpdatedAt
		photo["User"] = user

		data = append(data, photo)
	}

	c.JSON(http.StatusOK, data)

}

// Update godoc
// @Summary      Update an photo
// @Description  update an photo by ID
// @Tags         Photo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Photo ID"
// @Param        json  body  models.Photo true  "Photo"
// @Success      200  {object}  models.Photo
// @Security     Bearer
// @Router       /photos/{id}   [put]
func PhotoUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserId = userId
	Photo.ID = uint(photoId)

	err := db.Model(&Photo).Where("id = ?", photoId).Updates(models.Photo{Title: Photo.Title, Caption: Photo.Caption, PhotoUrl: Photo.PhotoUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.PhotoUrl,
		"user_id":    Photo.UserId,
		"updated_at": Photo.UpdatedAt,
	})
}

// Delete godoc
// @Summary      Delete an photo
// @Description  delete an photo by ID
// @Tags         Photo
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Photo ID"
// @Success      200  {string}  string
// @Security     Bearer
// @Router       /photos/{id}   [delete]
func PhotoDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userId := uint(userData["id"].(float64))

	Photo.UserId = userId
	Photo.ID = uint(photoId)

	err := db.Delete(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
