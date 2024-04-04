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
// @Summary      Create an socialMedia
// @Description  create and store an socialMedia
// @Tags         Social Media
// @Param        name formData string true "Social Media's Name"
// @Param        social_media_url formData string true "Social Media's Social Media URL"
// @Success      201  {object}  models.SocialMedia
// @Security 	bearerAuth
// @Router       /socialmedias  [post]
func SocialMediaCreate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = userID

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"created_at":       SocialMedia.CreatedAt,
	})
}

// GetSocialMediaByID godoc
// @Summary      Get social media by ID
// @Description  Retrieve social media by its ID
// @Tags         Social Media
// @Param        socialMediaId   path   int  true  "Social Media ID"
// @Success      200  {object}  models.SocialMedia
// @Security    BearerAuth
// @SecurityDefinitions   Bearer:
// @type  http
// @scheme bearer
// @bearerFormat Bearer <token>
// @name  Authorization
// @in    header
// @Router       /socialmedias/{socialMediaId} [get]
func GetSocialMediaByID(c *gin.Context) {
	db := database.GetDB()
	var data map[string]interface{}
	var socialMedia models.SocialMedia
	socialMediaID, err := strconv.Atoi(c.Param("socialMediaId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid social media ID",
		})
		return
	}

	if err := db.Preload("User").First(&socialMedia, socialMediaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "Social media not found",
		})
		return
	}

	user := make(map[string]interface{})

	user["id"] = socialMedia.User.ID
	user["email"] = socialMedia.User.Email
	user["username"] = socialMedia.User.Username

	data = map[string]interface{}{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaUrl,
		"user_id":          socialMedia.UserId,
		"user":             user,
	}

	c.JSON(http.StatusOK, data)
}

// Fetch godoc
// @Summary      Fetch socialMedias
// @Description  get socialMedias
// @Tags         Social Media
// @Success      200	{object}	[]models.SocialMedia
// @Security    BearerAuth
// @Router       /socialmedias  [get]
func SocialMediaList(c *gin.Context) {
	db := database.GetDB()
	var Socmed []models.SocialMedia

	var data []interface{}

	err := db.Preload("User").Find(&Socmed).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for i := range Socmed {
		sosmed := make(map[string]interface{})
		user := make(map[string]interface{})

		user["id"] = Socmed[i].User.ID
		user["username"] = Socmed[i].User.Username

		sosmed["id"] = Socmed[i].ID
		sosmed["name"] = Socmed[i].Name
		sosmed["social_media_url"] = Socmed[i].SocialMediaUrl
		sosmed["user_id"] = Socmed[i].UserId
		sosmed["created_at"] = Socmed[i].CreatedAt
		sosmed["updated_at"] = Socmed[i].UpdatedAt
		sosmed["user"] = user

		data = append(data, sosmed)
	}

	c.JSON(http.StatusOK, data)
}

// Update godoc
// @Summary      Update an socialMedia
// @Description  update an socialMedia by ID
// @Tags         Social Media
// @Param        socialMediaId   path      int  true  "SocialMedia ID"
// @Param        name formData string true "Social Media's Name"
// @Param        social_media_url formData string true "Social Media's Social Media URL"
// @Success      200  {string}  string
// @Security    BearerAuth
// @Router       /socialmedias/{socialMediaId} [put]
func SocialMediaUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}

	socialMediaID, _ := strconv.Atoi(c.Param("socialMediaId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserId = userId
	SocialMedia.ID = uint(socialMediaID)

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaID).Updates(models.SocialMedia{Name: SocialMedia.Name, SocialMediaUrl: SocialMedia.SocialMediaUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               SocialMedia.ID,
		"name":             SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id":          SocialMedia.UserId,
		"updated_at":       SocialMedia.UpdatedAt,
	})
}

// Delete godoc
// @Summary      Delete an socialMedia
// @Description  delete an socialMedia by ID
// @Tags         Social Media
// @Param        socialMediaId   path      int  true  "SocialMedia ID"
// @Success      200  {string}  string
// @Security    BearerAuth
// @Router       /socialmedias/{socialMediaId} [delete]
func SocialMediaDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	SocialMedia := models.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userId := uint(userData["id"].(float64))

	SocialMedia.ID = uint(socialMediaId)
	SocialMedia.UserId = userId

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaId).Delete(models.SocialMedia{}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
