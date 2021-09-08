package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/qa-api/pkg/questionary/config"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository/mongoDB"
	grpcserver "github.com/ismaeljpv/qa-api/pkg/questionary/server/grpc"
	httpserver "github.com/ismaeljpv/qa-api/pkg/questionary/server/http"
	"github.com/ismaeljpv/qa-api/pkg/questionary/service"
	grpctransport "github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc"
	pb "github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc/protobuff"
	httptransport "github.com/ismaeljpv/qa-api/pkg/questionary/transport/http"
	"google.golang.org/grpc"
)

const mongoDBURI = "MONGODB_URI"

func main() {
	var httpAddr = flag.String("http", ":8080", "HTTP listen address")
	var logger log.Logger
	var grpcAddr = ":50051"
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"service", "questionary",
		"time:", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	errs := make(chan error)
	connURI, confErr := config.GetConfig(mongoDBURI)
	if confErr != nil {
		panic(confErr)
	}

	level.Info(logger).Log("msg", fmt.Sprintf("Connection URI prepared -> %v", connURI))
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	flag.Parse()
	ctx := context.Background()

	repo, repoErr := mongoDB.NewRepository(ctx, logger, connURI)
	if repoErr != nil {
		panic(repoErr)
	}

	serv := service.NewService(repo, logger)
	httpEndpoints := httptransport.MakeEndpoints(serv)
	grpcEndpoints := grpctransport.MakeEndpoints(serv)

	grpcServer := grpcserver.NewGRPCServer(grpcEndpoints, logger)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("msg", fmt.Sprintf("GRPC Server started listening on port %v", grpcAddr))
		baseServer := grpc.NewServer()
		pb.RegisterQuestionaryServiceServer(baseServer, grpcServer)
		baseServer.Serve(grpcListener)
	}()

	go func() {
		handler := httpserver.NewHTTPServer(ctx, httpEndpoints)
		level.Info(logger).Log("msg", fmt.Sprintf("HTTP Server listening on port %v", *httpAddr))
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
	close(errs)
}
