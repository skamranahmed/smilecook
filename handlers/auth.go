package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/smilecook/config"
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
type AuthHandler struct {
	ctx         context.Context
	collection  *mongo.Collection
	userService service.UserService
}

type userSignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userSignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection, userService service.UserService) *AuthHandler {
	return &AuthHandler{
		ctx:         ctx,
		collection:  collection,
		userService: userService,
	}
}

func (handler *AuthHandler) SignUpHandler(c *gin.Context) {
	var request userSignUpRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameAlreadyExists, err := handler.userService.DoesUsernameAlreadyExist(request.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if usernameAlreadyExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	// hash the password before saving in db
	hashedPassword, err := handler.userService.HashPassword(request.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  request.Username,
		Password:  hashedPassword,
		IsAdmin:   false, // always false in this handler
	}

	err = handler.userService.Create(user)
	if err != nil {
		fmt.Printf("Unable to insert user: %+v in db, err: %v\n", user, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "signup successfull"})
	return
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var request userSignInRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find a user with the requested username
	user, err := handler.userService.FindOne(request.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// no user record found
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// verify the password
	err = handler.userService.VerifyPassword(request.Password, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtOutput := JWTOutput{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
	return
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := JWTOutput{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	c.JSON(http.StatusOK, jwtOutput)
	return

}
