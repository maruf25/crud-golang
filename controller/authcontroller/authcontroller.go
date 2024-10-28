package authcontroller

import (
	"EcommerceSederhana/config"
	"EcommerceSederhana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		panic(err)
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "email already taken",
		})
		return
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashPassword)

	if result := config.DB.Create(&user); result.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	c.JSON(http.StatusCreated, &user)
}

func Login(c *gin.Context) {
	// Ambil Request
	var user models.LoginUser

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	}
	// Cari apakah email ada pada database
	var existingUser models.User
	if err := config.DB.Where("email = ? ", user.Email).First(&existingUser).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Email or password is wrong",
			})
			return
		default:

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Request Error From Server",
			})
			return
		}
	}

	// Sesuaikan password apakah benar
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or password is wrong",
		})
		return
	}

	// Generate token
	expTime := time.Now().Add(time.Hour + 1)
	claims := &config.JWTClaim{
		Email: existingUser.Email,
		Role:  string(existingUser.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ecom-jwt",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// Algoritma Token
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenAlgo.SignedString(config.SecretKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "success login",
	})
}
