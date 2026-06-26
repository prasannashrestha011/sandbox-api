package main

import (
	"log"
	"main/internal/jobs/warmpool"

	"github.com/hibiken/asynq"
)

func main() {
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379", Password: "9843"}
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(warmpool.TypeWarmPool, warmpool.HandleSandboxProvision)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
