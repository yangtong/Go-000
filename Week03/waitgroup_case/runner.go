package waitgroup_case

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	Server *http.Server
	Wg     sync.WaitGroup
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

func (r *Runner) SignalHandler(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("starting signal goroutine")

	defer r.Wg.Done()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("signal ctx done")
			return

		case <-quit:
			cancel()
			return
		}
	}
}

func (r *Runner) StopHttpServer(ctx context.Context) {
	fmt.Println("starting stop goroutine")

	defer r.Wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("http stop ctx done")
			err := r.Server.Shutdown(ctx)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
}

func (r *Runner) StartHttpServer() {
	fmt.Println("starting http server")
	
	defer r.Wg.Done()

	err := r.Server.ListenAndServe()
	if err != nil {
		fmt.Println("http server done, err: ", err)
	}
}

func (r *Runner) StartBizServer(ctx context.Context) {
	fmt.Println("starting biz goroutine")
	
	defer r.Wg.Done()

	for {
		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("biz ctx done")
			return

		case <-ticker.C:
			fmt.Println("biz sleep 5s")
		}
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("hello world")))
}
