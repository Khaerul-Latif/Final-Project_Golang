package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var appJSON = "application/json"

// Register godoc
// @Summary      Create an user
// @Description  create and store an user
// @Tags         User
// @Param        email formData string true "User's Email"
// @Param        username formData string true "User's Username"
// @Param        age formData int true "User's Age"
// @Param        password formData string true "User's Password"
// @Param        profile_image_url formData string true "User's Profile Image URL"
// @Success      201  {object}   models.User
// @Router       /users/register [post]
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	user := models.User{}

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := c.ShouldBind(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	if err := db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": "Email already exists",
			})

			return
		}

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": "Username already exists",
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"age":               user.Age,
		"email":             user.Email,
		"id":                user.ID,
		"username":          user.Username,
		"profile_image_url": user.ProfileImageURL,
	})
}

// Login godoc
// @Summary      Show an user
// @Description  get an user by ID
// @Tags         User
// @Param        email formData string true "User's Email"
// @Param        password formData string true "User's Password"
// @Success      200  {object}  models.User
// @Router       /users/login	  [post]
func UserLogin(c *gin.Context) {
	db := database.GetDB()

	contentType := helpers.GetContentType(c)
	user := models.User{}

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := c.ShouldBind(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	originalPassword := user.Password
	if err := db.Where("email = ?", user.Email).First(&user).Take(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Not Found",
			"message": "User not found",
		})

		return
	}

	if isValid := helpers.CheckPasswordHash([]byte(user.Password), []byte(originalPassword)); !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})

		return
	}

	jwt := helpers.GenerateToken(user.ID, user.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
	})

}

// Update godoc
// @Summary      Update an user
// @Description  update an user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        email formData string true "User's Email"
// @Param        username formData string true "User's Username"
// @Param        age formData int true "User's Age"
// @Param        password formData string true "User's Password"
// @Param        profile_image_url formData string true "User's Profile Image URL"
// @Success      200  {string}  models.User
// @Security    BearerAuth
// @Router       /users [put]
func UserUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	user := models.User{}

	userid, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid user id",
		})
		return
	}

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := c.ShouldBind(&user); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	user.ID = userID

	err = db.Model(&user).Where("id = ?", userid).Updates(&user).First(&user).Error

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": "Email already exists",
			})

			return
		}

		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Conflict",
				"message": "Username already exists",
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                user.ID,
		"email":             user.Email,
		"username":          user.Username,
		"updated_at":        user.UpdatedAt,
		"age":               user.Age,
		"profile_image_url": user.ProfileImageURL,
	})

}

// Delete godoc
// @Summary      Delete an user
// @Description  delete an user by ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  []interface{}  "Your comment has been successfully deleted"
// @Security    BearerAuth
// @Router       /users					[delete]
func UserDelete(ctx *gin.Context) {
	userData, exists := ctx.Get("userData")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "User data not found in context",
		})
		return
	}

	claims, ok := userData.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to parse user data",
		})
		return
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Invalid user ID data type",
		})
		return
	}

	userID := uint(userIDFloat)
	// Mendapatkan data pengguna dari database
	var user models.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve user data",
		})
		return
	}

	// Menghapus pengguna dari database
	if err := database.GetDB().Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "Failed to delete user account",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
