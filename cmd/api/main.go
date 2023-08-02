package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"os/signal"
	"project.lemfi.net/internal/data"
	"sync"
	"syscall"
	"time"
)

type config struct {
	port int
}

type application struct {
	logger            *zerolog.Logger
	store             *data.Store
	verificationQueue chan *data.User
	transactionQueue  chan *data.Transaction
	transactionMutex  sync.Mutex
	cfg               config
}

func main() {

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	app := setup(wg, ctx)
	router := app.routes()

	addr := fmt.Sprintf(":%d", app.cfg.port)
	app.logger.Info().Msg(fmt.Sprintf("starting server on %s", addr))
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go server.ListenAndServe()

	waitForShutdown() // wait for shutdown signal
	app.logger.Info().Msg("Shutting down")

	if err := server.Shutdown(ctx); err != nil {
		app.logger.Error().Msg("Failed to shut down the server gracefully")
	}
	cancel()  // cancel ctx deep in runtime stack
	wg.Wait() // wait for all goroutines to exit

}

func setup(group *sync.WaitGroup, ctx context.Context) *application {
	var waitTimeSecs int
	var numWorkers int
	var cfg config

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr, TimeFormat: time.RFC3339,
	}).Level(zerolog.InfoLevel).With().Timestamp().Logger()

	store := data.NewStore()
	app := &application{
		store:  store,
		logger: &logger,
	}

	flag.IntVar(&numWorkers, "num-workers", 10, "Number of workers")
	flag.IntVar(&waitTimeSecs, "worker-wait-time", 30, "Worker wait time (secs)")
	flag.IntVar(&cfg.port, "port", 9000, "Server port")
	flag.Parse()

	app.cfg = cfg

	jobQueue := make(chan *data.User, 100)
	transactionQueue := make(chan *data.Transaction, 100)
	app.verificationQueue = jobQueue
	app.transactionQueue = transactionQueue

	for i := 0; i < numWorkers; i++ {
		group.Add(2)
		go app.verificationWorker(group, ctx, waitTimeSecs)
		go app.transactionWorker(group, ctx, waitTimeSecs)
	}
	return app

}

func waitForShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-ch
}
