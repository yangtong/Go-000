package main

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

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("hello world")))
}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)

	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}
	wg.Add(3)

	go func() {
		fmt.Println("starting signal goroutine")
		defer wg.Done()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return

			case <-quit:
				cancel()
				err := server.Shutdown(ctx)
				fmt.Println(err)
				return
			}
		}
	}()

	go func() {
		fmt.Println("starting http server")

		defer wg.Done()

		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("http server done, err: ", err)
		}
	}()

	go func() {
		fmt.Println("starting biz goroutine")

		defer wg.Done()

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
	}()

	wg.Wait()

}
