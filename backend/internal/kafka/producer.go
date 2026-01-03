package kafka

import (
	"context"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Event struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Version   string      `json:"version"`
	Data      interface{} `json:"data"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishArticleCreated(ctx context.Context, article interface{}) error {
	return p.publishEvent(ctx, "portfolio.articles", "article.created", article)
}

func (p *Producer) PublishArticleUpdated(ctx context.Context, article interface{}) error {
	return p.publishEvent(ctx, "portfolio.articles", "article.updated", article)
}

func (p *Producer) PublishArticleDeleted(ctx context.Context, articleID string) error {
	return p.publishEvent(ctx, "portfolio.articles", "article.deleted", map[string]string{"id": articleID})
}

func (p *Producer) PublishProjectCreated(ctx context.Context, project interface{}) error {
	return p.publishEvent(ctx, "portfolio.projects", "project.created", project)
}

func (p *Producer) PublishProjectUpdated(ctx context.Context, project interface{}) error {
	return p.publishEvent(ctx, "portfolio.projects", "project.updated", project)
}

func (p *Producer) PublishProjectDeleted(ctx context.Context, projectID string) error {
	return p.publishEvent(ctx, "portfolio.projects", "project.deleted", map[string]string{"id": projectID})
}

func (p *Producer) publishEvent(ctx context.Context, topic, eventType string, data interface{}) error {
	event := Event{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Source:    "backend",
		Version:   "1.0",
		Data:      data,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Value: eventData,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte(eventType)},
			{Key: "content-type", Value: []byte("application/json")},
		},
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

