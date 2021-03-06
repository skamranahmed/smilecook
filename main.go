package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skamranahmed/smilecook/config"
	"github.com/skamranahmed/smilecook/handlers"
	"github.com/skamranahmed/smilecook/repository"
	"github.com/skamranahmed/smilecook/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx            context.Context
	err            error
	recipesHandler *handlers.RecipesHandler
	authHandler    *handlers.AuthHandler
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of incoming requests",
	},
	[]string{"path"},
)

var totalHTTPMethods = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_methods_total",
		Help: "Number of requests per HTTP method",
	},
	[]string{"method"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests",
	},
	[]string{"path"},
)

func init() {
	// mongodb client setup
	ctx = context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	err = mongoClient.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("✅ Connected to MongoDB")

	recipesCollection := mongoClient.Database(config.MongoDatabaseName).Collection("recipes")
	usersCollection := mongoClient.Database(config.MongoDatabaseName).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURI,
		Password: config.RedisPassword,
		DB:       0,
	})

	redisConnectionStatus := redisClient.Ping(ctx)
	redisConnectionErr := redisConnectionStatus.Err()
	if redisConnectionErr != nil {
		log.Fatalf("❌ unable to connect to redis, error: %v", redisConnectionErr)
	}

	log.Println("✅ Connected to Redis - ", redisConnectionStatus)

	// register promethues metrics
	prometheus.Register(totalRequests)
	prometheus.Register(totalHTTPMethods)
	prometheus.Register(httpDuration)

	// instantiate the repo(s)
	userRepository := repository.NewUserRepository(ctx, usersCollection)
	recipeRepository := repository.NewRecipeRepository(ctx, recipesCollection)

	// instantiate the service(s)
	userService := service.NewUserService(userRepository)
	recipeService := service.NewRecipeService(recipeRepository)

	// instantiate the handler(s)
	recipesHandler = handlers.NewRecipesHandler(ctx, recipesCollection, redisClient, recipeService)
	authHandler = handlers.NewAuthHandler(ctx, usersCollection, userService)
}

// this is just a test route - no logic here
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": os.Getenv("API_VERSION")})
	return
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(c.Request.URL.Path))
		totalRequests.WithLabelValues(c.Request.URL.Path).Inc()
		totalHTTPMethods.WithLabelValues(c.Request.Method).Inc()
		c.Next()
		// this needs to be placed after c.Next() as we need to observe the time once the request has been served
		timer.ObserveDuration()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecretKey), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("auth", claims)
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.Default())

	// Prometheus middleware
	router.Use(PrometheusMiddleware())

	router.GET("/version", VersionHandler)
	router.GET("/prometheus", gin.WrapH(promhttp.Handler()))
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signup", authHandler.SignUpHandler)
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.CreateRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	addr := fmt.Sprintf(":%s", config.ServerPort)
	router.Run(addr)
}
