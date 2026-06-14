package main

import (
	"log"
	"main/internal/jobs"

	"github.com/hibiken/asynq"
)

func main() {
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379", Password: "9843"}
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypeSandboxCleanup, jobs.HandleSandboxCleanup)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
