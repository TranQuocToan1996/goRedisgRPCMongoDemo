package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/TranQuocToan1996/redislearn/config"
	"github.com/TranQuocToan1996/redislearn/controllers"
	"github.com/TranQuocToan1996/redislearn/routes"
	"github.com/TranQuocToan1996/redislearn/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      = gin.Default()
	ctx         = context.Background()
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// Connect to MongoDB
	mongoconnOpt := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconnOpt)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping().Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set("test", "Welcome to Golang with Redis and MongoDB",
		0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get("test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	server.SetTrustedProxies(nil)

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	log.Fatal(server.Run(":" + config.Port))
}