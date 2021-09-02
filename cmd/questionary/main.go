package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository/mock"
	"github.com/ismaeljpv/qa-api/pkg/questionary/server"
	"github.com/ismaeljpv/qa-api/pkg/questionary/service"
	"github.com/ismaeljpv/qa-api/pkg/questionary/transport"
)

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"service", "questionary",
		"time:", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	flag.Parse()
	ctx := context.Background()

	repo := mock.NewRepository([]domain.QuestionInfo{}, logger)
	serv := service.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(serv)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := server.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}