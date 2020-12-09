package errgroup_case

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Runner struct {
	Server      *http.Server
	DebugServer *http.Server
}

func NewRunner() *Runner {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)

	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}

	debugServer := &http.Server{
		Addr:         "127.0.0.1:8081",
		WriteTimeout: time.Second * 5,
		Handler:      http.DefaultServeMux,
	}

	return &Runner{
		Server:      server,
		DebugServer: debugServer,
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("hello world")))
}

func (r *Runner) SignalRoutine(ctx context.Context) error {
	fmt.Println("starting signal goroutine")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-quit:
			return errors.New("quit by signal")
		}
	}
}

func (r *Runner) HttpRoutine() error {
	fmt.Println("starting http server")
	err := r.Server.ListenAndServe()
	if err != nil {
		fmt.Println("http server done, err: ", err)
		return err
	}
	return nil
}

func (r *Runner) DebugHttpRoutine() error {
	fmt.Println("starting debug http server")
	err := r.DebugServer.ListenAndServe()
	if err != nil {
		fmt.Println("debug http server done, err: ", err)
		return err
	}
	return nil
}

func (r *Runner) BizRoutine(ctx context.Context) error {
	fmt.Println("starting biz goroutine")
	for {		
		select {
		case <-ctx.Done():
			fmt.Println("biz ctx done")
			return ctx.Err()
		case <-time.After(time.Second * 5):
			fmt.Println("biz sleep 5s")
		}
	}
}

func (r *Runner) StopRoutine(ctx context.Context) error {
	fmt.Println("starting stop http goroutine")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop http ctx done")
			err := r.Server.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func (r *Runner) StopDebugRoutine(ctx context.Context) error {
	fmt.Println("starting stop debug http goroutine")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop debug http ctx done")
			err := r.DebugServer.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		}
	}
}
