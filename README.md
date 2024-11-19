# SQS consumer on Golang
> Use AWS SDK for Go V2 and zerolog

## Description

Very simple consumer for proceed messages from queue there we can set:
 - batchSize (Max 10, default 1)
 - waitTimeSeconds (Default is 0s)
 - MessageAttributeNames (Default [])
 - messageSystemAttributeNames (Default [])
 - pollDelayInMs (default 0)
 - visibilityTimeout (default 0)
 - shouldDeleteMessages flag for cases there don't need to delete the message (Default true - delete message)

Default log level is `info`
## Usage example

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	consumer "github.com/1xtr/go-sqs-consumer"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog"
)

func main() {
	log := consumer.GetLogger("main")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Error().Err(err).Msgf("error load sdk config: %v", err)
		return
	}
	opts := consumer.Options{
		QueueUrl:              "https://sqs.eu-central-1.amazonaws.com/1234567890/development",
		SqsClient:             sqs.NewFromConfig(cfg),
		HandleMessage:         handler,
		ShouldDeleteMessages:  aws.TrueTernary,
		WaitTimeSeconds:       1,
		MessageAttributeNames: []string{"All"},
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{
			types.MessageSystemAttributeNameAll,
		},
	}
	c := consumer.New(opts)
	c.Start()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Debug().Msg("Gracefully shutting down...")
	log.Debug().Msg("Running cleanup tasks...")

	c.Stop()
	log.Debug().Msg("App was successful shutdown.")
}

func handler(ctx context.Context, msg *types.Message) error {
	log := zerolog.Ctx(ctx).With().Str("MessageId", *msg.MessageId).
		Str("component", "handler").Logger()

	log.Info().Interface("message", msg).Msgf("message received")
	time.Sleep(5 * time.Second)
	if strings.Contains(*msg.Body, "error") {
		return fmt.Errorf("an error occured")
	}
	return nil
}
```
---

### Tests

For local tests we can set environment variable `LOG_TO_FILE=true` for save execution log to `./logs` folder.
