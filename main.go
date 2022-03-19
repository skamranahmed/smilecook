package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/skamranahmed/smilecook/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	recipes         []Recipe
	ctx             context.Context
	err             error
	mongoClient     *mongo.Client
	collection      *mongo.Collection
	collectionUsers *mongo.Collection
	recipesHandler  *handlers.RecipesHandler
	authHandler     *handlers.AuthHandler
)

func init() {
	// mongodb client setup
	ctx = context.Background()
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Connected to MongoDB")

	collection = mongoClient.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collectionUsers = mongoClient.Database(os.Getenv("MONGO_DATABASE")).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	redisConnectionStatus := redisClient.Ping(ctx)
	redisConnectionErr := redisConnectionStatus.Err()
	if redisConnectionErr != nil {
		log.Fatalf("❌ unable to connect to redis, error: %v", redisConnectionErr)
	}

	log.Println("✅ Connected to Redis - ", redisConnectionStatus)

	// instantiate the handler(s)
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
}

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")

// 	listOfRecipes := make([]Recipe, 0)

// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}

// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}

// 	c.JSON(http.StatusOK, listOfRecipes)
// 	return
// }

func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("API_VERSION")})
	return
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
		// 	c.AbortWithStatus(401)
		// }

		tokenValue := c.GetHeader("Authorization")

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.Default())

	router.GET("/version", VersionHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signup", authHandler.SignUpHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	router.Run()
}
