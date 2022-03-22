package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/skamranahmed/smilecook/models"
	"github.com/skamranahmed/smilecook/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RecipesHandler struct {
	ctx           context.Context
	collection    *mongo.Collection
	redisClient   *redis.Client
	recipeService service.RecipeService
}

// NewRecipesHandler: used to create a new instance from the RecipesHanlder struct
func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client, recipeService service.RecipeService) *RecipesHandler {
	return &RecipesHandler{
		ctx:           ctx,
		collection:    collection,
		redisClient:   redisClient,
		recipeService: recipeService,
	}
}

// CreateRecipeHandler: inserts a new recipe
func (handler *RecipesHandler) CreateRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	err = handler.recipeService.Create(&recipe)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	log.Println("deleting data from redis")
	handler.redisClient.Del(handler.ctx, "recipes")

	c.JSON(http.StatusOK, recipe)
	return
}

// ListRecipesHandler: fetches a list of recipes
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(handler.ctx, "recipes").Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("value not found in redis, hitting mongo db now")
		} else {
			log.Printf("error in retrieving value from redis, err: %v, hitting mongo db now\n", err)
		}

		// fetch all recipes from mongo db
		recipes, err := handler.recipeService.FetchAll()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// save the data in redis
		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(handler.ctx, "recipes", string(data), 0)

		c.JSON(http.StatusOK, recipes)
		return
	}

	log.Println("request to redis")
	recipes := make([]*models.Recipe, 0)
	json.Unmarshal([]byte(val), &recipes)
	c.JSON(http.StatusOK, recipes)
	return
}

func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe models.Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordExists, err := handler.recipeService.Update(objectID, &recipe)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err == nil && !recordExists {
		// this means no recipe record was found for the requested id, but the operation succeeded without any error
		errMsg := fmt.Sprintf("no recipe found with id: %s", id)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	log.Println("deleting data from redis")
	handler.redisClient.Del(handler.ctx, "recipes")

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
	return
}

func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = handler.collection.DeleteOne(handler.ctx, bson.M{"_id": objectID})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

func (handler *RecipesHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{"_id": objectID})

	var recipe models.Recipe
	err = cur.Decode(&recipe)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
	return
}
