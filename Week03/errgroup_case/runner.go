package errgroup_case

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Runner struct {
	Server *http.Server
}

func NewRunner() *Runner {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)

	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}

	return &Runner{
		Server: server,
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

func (r *Runner) BizRoutine(ctx context.Context) error {
	fmt.Println("starting biz goroutine")
	for {
		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("biz ctx done")
			return ctx.Err()
		case <-ticker.C:
			fmt.Println("biz sleep 5s")
			return nil
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
