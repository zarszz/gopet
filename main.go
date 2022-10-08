package main

import (
	"context"
	"fmt"
	"go-grpc/config"
	"go-grpc/controllers"
	"go-grpc/gapi"
	"go-grpc/pb"
	"go-grpc/routes"
	"go-grpc/services"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	postCollection *mongo.Collection
	postService    services.PostService

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
	postCollection = mongoClient.Database("go_grpc").Collection("posts")
	userService = services.NewUserServiceImpl(authCollection, ctx)

	authService = services.NewAuthService(authCollection, ctx)
	authController = controllers.NewAuthController(authService, userService, temp)
	authRouteController = routes.NewAuthRouteController(authController)

	userController = controllers.NewUserController(userService)
	userRouteController = routes.NewRouteUserController(userController)

	postService = services.NewPostService(postCollection, ctx)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)
	// startGinServer(config)
	startGrpcServer(config, authCollection, postCollection)
}

func startGinServer(config config.Config) {
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

func startGrpcServer(config config.Config, aCollection *mongo.Collection, pCollection *mongo.Collection) {
	authServer, err := gapi.NewGrpcAuthServer(
		config, authService, userService, aCollection, temp,
	)
	if err != nil {
		log.Fatalf("cannot create auth grpc server : %v", err)
	}

	userServer, err := gapi.NewGrpcUserServer(
		config, userService,
	)
	if err != nil {
		log.Fatalf("cannot create user grpc server : %v", err)
	}

	postServer, err := gapi.NewGrpcPostServer(pCollection, postService)
	if err != nil {
		log.Fatalf("cannot create post grpc server : %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterUserServiceServer(grpcServer, userServer)
	pb.RegisterPostServiceServer(grpcServer, postServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("cannot create grpc server : %v", err)
	}

	log.Printf("start gRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot create grpc server : %v", err)
	}
}
