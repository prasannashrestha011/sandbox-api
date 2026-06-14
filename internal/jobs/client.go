package jobs

import "github.com/hibiken/asynq"

func InitAsynq() *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     "127.0.0.1:6379",
		Password: "9843",
	})
	return client
}
