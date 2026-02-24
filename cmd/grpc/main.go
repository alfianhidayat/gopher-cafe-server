package main

import (
	"gopher-cafe/internal/entity/coffeeshop"
	"gopher-cafe/internal/worker"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	appCfg "gopher-cafe/config"
	handler "gopher-cafe/internal/handler/grpc/coffeeshop"
	usecase "gopher-cafe/internal/usecase/coffeeshop"

	pb "github.com/rexyajaib/gopher-cafe/pkg/gen/go/v1"

	"github.com/ajaibid/coin-common-golang/config"
	"github.com/ajaibid/coin-common-golang/logger"

	_ "net/http/pprof"
)

func main() {
	var (
		equipPoolManager *worker.EquipPoolManager
		grpcServer       *grpc.Server
	)

	shutdown := func() {
		logger.Info("Begin Shutting down gracefully...")
		if grpcServer != nil {
			logger.Info("Shutting down grpc server...")
			grpcServer.Stop()
		}
		if equipPoolManager != nil {
			logger.Info("Shutting down manager...")
			equipPoolManager.StopAll()
		}
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Recovered from panic: %v, stack trace: %s", r, debug.Stack())
			shutdown()
			os.Exit(1)
		}
	}()

	// Create a TCP Listener on a specific port
	var cfg appCfg.Config
	err := config.LoadConfig(&cfg, "config/.env")
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Grpc.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.NewLogger(logger.LoggerConf{
		LogLevel:     cfg.Logger.LogLevel,
		LogFormatter: cfg.Logger.LogFormatter,
	})

	equipWorkers := worker.EquipmentWorkers

	equipPoolManager = worker.NewEquipPoolManager(uint8(len(equipWorkers)))
	for k, v := range equipWorkers {
		equipPoolManager.Register(k, v)
	}
	equipPoolManager.StartAll()

	metrics := coffeeshop.NewOrderMetrics()

	// Initialize the Layers
	coffeeUsecase := usecase.NewCoffeeshopUsecase(equipPoolManager, metrics)
	coffeeHandler := handler.NewCoffeeshopGrpcHandler(coffeeUsecase)

	// Create the gRPC Server instance
	grpcServer = grpc.NewServer(grpc.UnaryInterceptor(TimeoutMiddleware()))

	// Register the Service (The "Route Definition")
	// This tells the gRPC server to route incoming GopherCafe calls to our handler.
	pb.RegisterGopherCafeServiceServer(grpcServer, coffeeHandler)

	// Optional: Enable reflection.
	// This allows tools like Postman or 'evans' to "see" your endpoints automatically.
	reflection.Register(grpcServer)

	// Start Serving
	log.Printf("Coffee Shop Simulation Server is running on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	logger.Info("Received shutdown signal, shutting down...")
	shutdown()
}
