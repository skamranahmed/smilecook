package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/skamranahmed/smilecook/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	recipes        []Recipe
	ctx            context.Context
	err            error
	mongoClient    *mongo.Client
	collection     *mongo.Collection
	recipesHandler *handlers.RecipesHandler
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	log.Println("✅ Connected to Redis, PING:", status)

	// instantiate the recipes handler
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
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

func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)
	router.Run()
}
