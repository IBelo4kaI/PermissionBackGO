package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	grpcpb "permisson/internal/adapter/grpc"
	"permisson/internal/database"
	repo "permisson/internal/database/sqlc"
	"permisson/internal/pkg/env"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	env := env.Env{}
	env.Init()

	sessionTTL, err := env.GetSessionTTL(4 * time.Hour)
	if err != nil {
		panic(err)
	}

	cfg := config{
		addr: env.GetAddr(),
		db: dbConfig{
			dsn: env.GetDbString(),
		},
		prefix:       "time",
		appEnv:       env.GetAppEnv(),
		cookieDomain: env.GetCookieDomain(),
		sessionTTL:   sessionTTL,
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      "15:04:05",
	})

	db, err := sql.Open("mysql", cfg.db.dsn)
	if err != nil {
		fmt.Println("error opening database:", err)
		panic("error con database")
	}
	defer db.Close()

	if err := database.RunMigrations(context.Background(), db); err != nil {
		panic("migrations failed: " + err.Error())
	}

	// // gRPC-сервер (как в Python-версии — отдельный порт, запускается параллельно с HTTP).
	// go func() {
	// 	if err := runGRPCServer(repo.New(db)); err != nil {
	// 		logger.Error("gRPC server failed", "error", err)
	// 	}
	// }()

	app := application{
		config: cfg,
		db:     db,
		logger: logger,
	}

	// app.run(app.mount())

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return runGRPCServer(repo.New(db))
	})

	g.Go(func() error {
		return app.run(app.mount())
	})

	if err := g.Wait(); err != nil {
		logger.Fatal(err)
	}
}

// runGRPCServer запускает PermissionService и UserService на localhost:8382
// (тот же порт, что в Python-версии).
func runGRPCServer(queries repo.Querier) error {
	log.Info("starting gRPC server")
	lis, err := net.Listen("tcp", ":8383")
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}

	server := grpc.NewServer()
	grpcSvc := grpcpb.NewServers(queries)
	grpcpb.RegisterPermissionServiceServer(server, grpcSvc)
	grpcpb.RegisterUserServiceServer(server, grpcSvc)

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("gRPC serve: %w", err)
	}
	return nil
}
