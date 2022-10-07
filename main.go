package main

import (
	"context"
	"fmt"
	"go-grpc/config"
	"go-grpc/controllers"
	"go-grpc/routes"
	"go-grpc/services"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client

	userService         services.UserService
	userController      controllers.UserController
	userRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	authController      controllers.AuthController
	authRouteController routes.AuthRouteController

	temp *template.Template
)

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variable", err)
	}

	ctx = context.TODO()

	// connect to mongoDB
	mongoConnection := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConnection)
	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected")

	// connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisUri,
		Password: config.RedisPassword,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, "test", "welcome to jungle", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis successfully connected")

	// init dep
	authCollection := mongoClient.Database("go_grpc").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)

	authService = services.NewAuthService(authCollection, ctx)
	authController = controllers.NewAuthController(authService, userService, temp)
	authRouteController = routes.NewAuthRouteController(authController)

	userController = controllers.NewUserController(userService)
	userRouteController = routes.NewRouteUserController(userController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()
	if err == redis.Nil {
		fmt.Println("key: test does not exists")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	authRouteController.AuthRoute(router, userService)
	userRouteController.UserRoute(router, userService)

	log.Fatal(server.Run(":" + config.Port))
}
