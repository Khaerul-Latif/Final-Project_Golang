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

// Store godoc
// @Summary      Create an comment
// @Description  create and store an comment
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        message formData string true "Comments's Message"
// @Param        photo_id formData int true "Comments's Photo ID"
// @Success      201  {object}  models.Comment
// @Security     Bearer
// @Router       /comments      [post]
func CommentCreate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserId = userID

	err := db.Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoId,
		"user_id":    Comment.UserId,
		"created_at": Comment.CreatedAt,
	})
}

// Fetch godoc
// @Summary      Fetch comments
// @Description  get comments
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Success      200	{object}	[]models.Comment
// @Security     Bearer
// @Router       /comments      [get]
func CommentList(c *gin.Context) {
	db := database.GetDB()
	Comments := []models.Comment{}
	var data []interface{}

	err := db.Model(&models.Comment{}).Preload("User").Preload("Photo").Find(&Comments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for i := range Comments {
		photo := make(map[string]interface{})
		user := make(map[string]interface{})

		user["id"] = Comments[i].User.ID
		user["email"] = Comments[i].User.Email
		user["username"] = Comments[i].User.Username

		photo["id"] = Comments[i].Photo.ID
		photo["title"] = Comments[i].Photo.Title
		photo["caption"] = Comments[i].Photo.Caption
		photo["photo_url"] = Comments[i].Photo.PhotoUrl
		photo["user_id"] = Comments[i].Photo.UserId

		data = append(data, gin.H{
			"id":         Comments[i].ID,
			"message":    Comments[i].Message,
			"photo_id":   Comments[i].PhotoId,
			"user_id":    Comments[i].UserId,
			"created_at": Comments[i].CreatedAt,
			"updated_at": Comments[i].UpdatedAt,
			"User":       user,
			"Photo":      photo,
		})
	}

	c.JSON(http.StatusOK, data)
}

// CommentByID godoc
// @Summary      Get a comment by ID
// @Description  Retrieve a comment by its ID
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {object}  models.Comment
// @Security     Bearer
// @Router       /comments/{commentId} [get]
func CommentByID(c *gin.Context) {
	db := database.GetDB()
	commentID := c.Param("commentId")
	var comment models.Comment
	var data map[string]interface{}

	if err := db.Preload("User").Preload("Photo").First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	user := make(map[string]interface{})
	photo := make(map[string]interface{})

	user["id"] = comment.User.ID
	user["email"] = comment.User.Email
	user["username"] = comment.User.Username

	photo["id"] = comment.Photo.ID
	photo["title"] = comment.Photo.Title
	photo["caption"] = comment.Photo.Caption
	photo["photo_url"] = comment.Photo.PhotoUrl
	photo["user_id"] = comment.Photo.UserId

	data = gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoId,
		"user_id":    comment.UserId,
		"created_at": comment.CreatedAt,
		"updated_at": comment.UpdatedAt,
		"User":       user,
		"Photo":      photo,
	}

	c.JSON(http.StatusOK, data)
}

// Update godoc
// @Summary      Update an comment
// @Description  update an comment by ID
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Param        message formData string true "Comments's Message"
// @Success      200  {object}  models.Photo
// @Security     Bearer
// @Router       /comments/{commentId} [put]
func CommentUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}
	photo := []models.Comment{}
	var data []interface{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.ID = uint(commentId)
	Comment.UserId = uint(userId)

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).First(&Comment).Error
	res := db.Model(&Comment).Preload("Photo").Where("id = ?", commentId).First(&photo).Error

	if res != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": res.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for i := range photo {
		photos := make(map[string]interface{})

		photos["id"] = photo[i].Photo.ID
		photos["title"] = photo[i].Photo.Title
		photos["caption"] = photo[i].Photo.Caption
		photos["photo_url"] = photo[i].Photo.PhotoUrl
		photos["user_id"] = photo[i].Photo.UserId
		photos["updated_at"] = photo[i].Photo.UpdatedAt

		data = append(data, photos)

	}

	c.JSON(http.StatusOK, data)
}

// Delete godoc
// @Summary      Delete an comment
// @Description  delete an comment by ID
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Comment ID"
// @Success      200  {string}  string
// @Security     Bearer
// @Router       /comments/{commentId} [delete]
func CommentDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	Comment := models.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userId := uint(userData["id"].(float64))

	Comment.ID = uint(commentId)
	Comment.UserId = userId

	err := db.Model(&Comment).Where("id = ?", commentId).Delete(models.Comment{}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
