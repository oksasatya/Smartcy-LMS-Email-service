package main

import (
	"github.com/ghssni/Smartcy-LMS/Email-Service/database"
	"github.com/ghssni/Smartcy-LMS/Email-Service/helper"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/repository"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/service"
	pb "github.com/ghssni/Smartcy-LMS/Email-Service/pb/proto"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

var db *gorm.DB

func main() {
	// Setup logger
	helper.SetupLogger()
	//.env
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// Config Db
	var err error
	db, err = database.InitDB()
	if err != nil {
		logrus.Fatalf("Error connecting to database: %v", err)
	}

	runGrpcServer()
}

func runGrpcServer() {
	grpcServer := grpc.NewServer()

	// initialize all services
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost"
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052"
	}
	address := grpcHost + ":" + grpcPort

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(helper.LogrusLoggerUnaryInterceptor),
	)

	//register service
	emailRepo := repository.NewEmailsRepository(db)
	emailLogRepo := repository.NewEmailsLogRepository(db)
	emailService := service.NewEmailService(emailRepo, emailLogRepo)

	// Register gRPC services
	pb.RegisterEmailServiceServer(grpcServer, emailService)

	// Start gRPC server
	log.Printf("gRPC server listening on %s", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
