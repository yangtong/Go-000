package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(string("hello world")))
}

func main() {

	group, ctx := errgroup.WithContext(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)

	server := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}

	group.Go(func() error {
		fmt.Println("starting signal goroutine")

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				err := server.Shutdown(ctx)
				if err != nil {
					return err
				}
				return errors.New("quit")
			}
		}
	})

	group.Go(func() error {
		fmt.Println("starting http server")
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("http server done, err: ", err)
			return err
		}
		return nil
	})

	group.Go(func() error {
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
	})

	err := group.Wait()
	if err == nil {
		fmt.Println("all finished ok")
	} else {
		fmt.Printf("finished error:%v\n", err)
	}

}
