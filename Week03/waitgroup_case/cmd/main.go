package main

import (
	"context"

	"github.com/yangtong/Go-000/Week03/waitgroup_case"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	runner := waitgroup_case.NewRunner()

	runner.Wg.Add(4)

	go runner.SignalHandler(ctx, cancel)
	go runner.StopHttpServer(ctx)
	go runner.StartHttpServer()
	go runner.StartBizServer(ctx)

	runner.Wg.Wait()
}
