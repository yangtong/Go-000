package main

import (
	"context"
	"fmt"
	"learn/Go-000/Week03/errgroup_case"

	"golang.org/x/sync/errgroup"
)

func main() {

	group, ctx := errgroup.WithContext(context.Background())

	runner := errgroup_case.NewRunner()

	group.Go(func() error {
		return runner.SignalRoutine(ctx)
	})

	group.Go(func() error {
		return runner.HttpRoutine()
	})

	group.Go(func() error {
		return runner.DebugHttpRoutine()
	})

	group.Go(func() error {
		return runner.BizRoutine(ctx)
	})

	group.Go(func() error {
		return runner.StopRoutine(ctx)
	})

	group.Go(func() error {
		return runner.StopDebugRoutine(ctx)
	})

	err := group.Wait()
	if err == nil {
		fmt.Println("all finished ok")
	} else {
		fmt.Printf("finished error:%v\n", err)
	}

}
