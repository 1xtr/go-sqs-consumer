package consumer

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rs/zerolog"
)

type (
	Consumer struct {
		sqsClient                   *sqs.Client
		queueUrl                    string
		handler                     func(c context.Context, m *types.Message) error
		stopSignal                  chan os.Signal
		messagesChannel             chan types.Message
		batchSize                   int
		pollDelayInMs               time.Duration
		visibilityTimeout           int
		waitTimeSeconds             int
		MessageAttributeNames       []string
		messageSystemAttributeNames []types.MessageSystemAttributeName
		shouldDeleteMessages        bool
		logger                      zerolog.Logger
	}
	Options struct {
		QueueUrl                    string
		SqsClient                   *sqs.Client
		BatchSize                   int
		PollDelayInMs               int
		VisibilityTimeout           int
		WaitTimeSeconds             int
		MessageAttributeNames       []string
		MessageSystemAttributeNames []types.MessageSystemAttributeName
		HandleMessage               func(c context.Context, m *types.Message) error
		ShouldDeleteMessages        aws.Ternary
	}
)

var (
	defaultRegion = "eu-central-1"
	region        = os.Getenv("AWS_REGION")
)

func New(o Options) *Consumer {
	c := Consumer{
		queueUrl:                    o.QueueUrl,
		sqsClient:                   o.SqsClient,
		handler:                     o.HandleMessage,
		stopSignal:                  make(chan os.Signal, 1),
		messagesChannel:             make(chan types.Message),
		batchSize:                   o.BatchSize,
		pollDelayInMs:               time.Duration(o.PollDelayInMs) * time.Millisecond,
		visibilityTimeout:           o.VisibilityTimeout,
		waitTimeSeconds:             o.WaitTimeSeconds,
		MessageAttributeNames:       o.MessageAttributeNames,
		messageSystemAttributeNames: o.MessageSystemAttributeNames,
		shouldDeleteMessages:        true,
		logger:                      GetLogger("Consumer"),
	}
	// If SQS Client not set, use default
	if c.sqsClient == nil {
		if region == "" {
			region = defaultRegion
		}
		c.sqsClient = sqs.New(
			sqs.Options{
				Region:           region,
				RetryMaxAttempts: 3,
				RetryMode:        "adaptive",
			},
		)
	}

	// If batch size not set use default 1
	if c.batchSize == 0 {
		c.batchSize = 1
	}
	// AWS batch size limit of 10,
	// https//docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/quotas-messages.html
	if c.batchSize > 10 {
		c.batchSize = 10
	}

	if o.ShouldDeleteMessages != aws.UnknownTernary {
		c.shouldDeleteMessages = o.ShouldDeleteMessages.Bool()
	}

	return &c
}

func (c *Consumer) Start() {
	log := GetLogger("Start")

	log.Debug().
		Str("queue", c.queueUrl).
		Msgf("consumer starting")

	signal.Notify(c.stopSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start message processing in a single goroutine
	go c.processMessages()

	// Poll messages in a loop
	c.pollMessages()

	// Wait for stop signal
	<-c.stopSignal
	c.logger.Debug().Msgf("Shutdown signal received. Stopping...")
	close(c.messagesChannel)
}

// waitForProcessing ensures that all messages in the current batch are processed before fetching new messages.
func (c *Consumer) waitForProcessing() {
	// Wait for the message channel to drain
	for {
		// If the message channel is empty, wait for pollDelay before re-checking
		if len(c.messagesChannel) == 0 {
			time.Sleep(c.pollDelayInMs)
			return
		}
		// If the message channel is not empty, wait for it to drain
		time.Sleep(100 * time.Millisecond) // Small delay to re-check channel status
	}
}

func (c *Consumer) pollMessages() {
	log := GetLogger("pollMessages")
	for {
		select {
		case <-c.stopSignal:
			log.Debug().Msgf("stop signal received, shutting down message receiver")
			close(c.messagesChannel)
			return
		default:
			result, err := c.sqsClient.ReceiveMessage(
				context.TODO(),
				&sqs.ReceiveMessageInput{
					MessageSystemAttributeNames: c.messageSystemAttributeNames,
					MessageAttributeNames:       c.MessageAttributeNames,
					QueueUrl:                    aws.String(c.queueUrl),
					MaxNumberOfMessages:         int32(c.batchSize),
					VisibilityTimeout:           int32(c.visibilityTimeout),
					WaitTimeSeconds:             int32(c.waitTimeSeconds),
				},
			)

			if err != nil {
				log.Error().Err(err).Caller().Msgf("error receive messages: %v", err)
				return
			}
			log.Debug().Interface("result", result).Msgf("pollMessages.result")
			if len(result.Messages) > 0 {
				for _, message := range result.Messages {
					c.messagesChannel <- message
				}
			}

			c.waitForProcessing()
		}
	}
}

func (c *Consumer) processMessages() {
	ctx := Logger.WithContext(context.Background())
	log := GetLogger("processMessages")
	for msg := range c.messagesChannel {
		err := c.handler(ctx, &msg)
		if err != nil {
			log.Error().Err(err).
				Interface("message", msg).
				Msgf("Error processing message: %v\n", err)
			continue
		}

		// Delete the message from SQS after successful processing
		if c.shouldDeleteMessages {
			go c.deleteMessage(ctx, &msg)
		}
	}
}

// Stop gracefully shuts down the consumer.
func (c *Consumer) Stop() {
	close(c.stopSignal)
}

func (c *Consumer) deleteMessage(ctx context.Context, msg *types.Message) {
	log := zerolog.Ctx(ctx).With().Str("MessageId", *msg.MessageId).
		Str("component", "deleteMessage").Logger()
	// Delete the message from SQS after successful processing
	if c.shouldDeleteMessages {
		log.Debug().Msgf("deleting message %s", *msg.MessageId)
		_, err := c.sqsClient.DeleteMessage(
			context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			},
		)
		if err != nil {
			log.Error().Err(err).Msgf("error deleting message %v", err)
		}
	}
}
