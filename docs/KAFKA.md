# Kafka Integration Documentation

## Genel Bakış

Portfolio Dashboard, event-driven architecture için Apache Kafka kullanır. Kafka, servisler arası asenkron iletişim, event sourcing ve microservices orchestration için kullanılır.

## Kullanım Senaryoları

### 1. Event-Driven Architecture

Sistem içindeki önemli olaylar Kafka topic'lerine publish edilir:

- **Article Events**: Makale oluşturma, güncelleme, silme
- **Project Events**: Proje ekleme, güncelleme, silme
- **Portfolio Events**: Portfolio bilgileri güncelleme
- **User Events**: Kullanıcı kaydı, giriş (auth service'den)

### 2. Async Processing

Uzun süren işlemler asenkron olarak işlenir:

- Email gönderimi
- Analytics event'leri
- Search index güncelleme
- Cache invalidation
- Image processing

### 3. Service Decoupling

Servisler birbirinden bağımsız çalışır:

- Backend → Kafka → Email Service
- Backend → Kafka → Analytics Service
- Backend → Kafka → Search Service

## Topic Yapısı

### Topics

#### portfolio.articles
**Partitions**: 3  
**Replication Factor**: 1 (local), 3 (production)  
**Retention**: 7 days

**Events**:
- `article.created`
- `article.updated`
- `article.deleted`
- `article.published`
- `article.unpublished`

#### portfolio.projects
**Partitions**: 3  
**Replication Factor**: 1 (local), 3 (production)  
**Retention**: 7 days

**Events**:
- `project.created`
- `project.updated`
- `project.deleted`
- `project.featured`

#### portfolio.analytics
**Partitions**: 6  
**Replication Factor**: 1 (local), 3 (production)  
**Retention**: 30 days

**Events**:
- `page.view`
- `article.view`
- `project.view`
- `user.action`

#### portfolio.notifications
**Partitions**: 3  
**Replication Factor**: 1 (local), 3 (production)  
**Retention**: 1 day

**Events**:
- `email.send`
- `notification.create`

## Event Schema

### Article Created Event

```json
{
  "event_id": "uuid",
  "event_type": "article.created",
  "timestamp": "2024-01-01T00:00:00Z",
  "source": "backend",
  "version": "1.0",
  "data": {
    "article_id": "uuid",
    "title": "Article Title",
    "slug": "article-slug",
    "author_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Project Created Event

```json
{
  "event_id": "uuid",
  "event_type": "project.created",
  "timestamp": "2024-01-01T00:00:00Z",
  "source": "backend",
  "version": "1.0",
  "data": {
    "project_id": "uuid",
    "name": "Project Name",
    "github_url": "https://github.com/user/repo",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## Producer Implementation

### Go Producer

```go
// internal/kafka/producer.go
package kafka

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
)

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

func (p *Producer) PublishArticleCreated(ctx context.Context, article *model.Article) error {
    event := Event{
        EventID:    uuid.New().String(),
        EventType:  "article.created",
        Timestamp:  time.Now(),
        Source:     "backend",
        Version:    "1.0",
        Data:       article,
    }
    
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    msg := kafka.Message{
        Topic: "portfolio.articles",
        Key:   []byte(article.ID),
        Value: data,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte("article.created")},
            {Key: "content-type", Value: []byte("application/json")},
        },
    }
    
    return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
    return p.writer.Close()
}
```

### Usage in Service

```go
// internal/service/article.go
func (s *articleService) CreateArticle(ctx context.Context, article *model.Article) error {
    // Save to database
    if err := s.repo.Create(ctx, article); err != nil {
        return err
    }
    
    // Publish event
    if err := s.kafkaProducer.PublishArticleCreated(ctx, article); err != nil {
        // Log error but don't fail the operation
        s.logger.Error("Failed to publish article.created event", zap.Error(err))
    }
    
    return nil
}
```

## Consumer Implementation

### Email Service Consumer

```go
// email-service/internal/kafka/consumer.go
package kafka

import (
    "context"
    "encoding/json"
    "github.com/segmentio/kafka-go"
)

type Consumer struct {
    reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:  brokers,
            Topic:    topic,
            GroupID:  groupID,
            MinBytes: 10e3,
            MaxBytes: 10e6,
        }),
    }
}

func (c *Consumer) ConsumeArticles(ctx context.Context, handler func(Event) error) error {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        var event Event
        if err := json.Unmarshal(msg.Value, &event); err != nil {
            continue
        }
        
        switch event.EventType {
        case "article.created":
            // Send welcome email to subscribers
            if err := handler(event); err != nil {
                // Log error, continue processing
                log.Printf("Error handling event: %v", err)
            }
        }
    }
}
```

### Analytics Service Consumer

```go
// analytics-service/internal/kafka/consumer.go
func (c *Consumer) ConsumeAnalytics(ctx context.Context) error {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        var event Event
        json.Unmarshal(msg.Value, &event)
        
        // Store in analytics database
        c.analyticsRepo.StoreEvent(ctx, event)
    }
}
```

## Consumer Groups

### Email Service
- **Group ID**: `email-service`
- **Topics**: `portfolio.articles`, `portfolio.notifications`
- **Purpose**: Email gönderimi

### Analytics Service
- **Group ID**: `analytics-service`
- **Topics**: `portfolio.analytics`
- **Purpose**: Analytics toplama

### Cache Service
- **Group ID**: `cache-service`
- **Topics**: `portfolio.articles`, `portfolio.projects`
- **Purpose**: Cache invalidation

### Search Service
- **Group ID**: `search-service`
- **Topics**: `portfolio.articles`, `portfolio.projects`
- **Purpose**: Search index güncelleme

## Error Handling

### Retry Strategy

```go
func (c *Consumer) ConsumeWithRetry(ctx context.Context, maxRetries int) error {
    for {
        msg, err := c.reader.ReadMessage(ctx)
        if err != nil {
            return err
        }
        
        var retries int
        for retries < maxRetries {
            if err := c.processMessage(ctx, msg); err == nil {
                break
            }
            retries++
            time.Sleep(time.Second * time.Duration(retries))
        }
        
        if retries >= maxRetries {
            // Send to dead letter queue
            c.sendToDLQ(ctx, msg)
        }
    }
}
```

### Dead Letter Queue

Failed message'lar için DLQ topic'i:

- `portfolio.articles.dlq`
- `portfolio.projects.dlq`
- `portfolio.analytics.dlq`

## Monitoring

### Metrics

- **Message Rate**: Topic başına message rate
- **Consumer Lag**: Consumer group lag
- **Error Rate**: Failed message rate
- **Throughput**: Messages per second

### Kafka Manager / Kafdrop

Kafka cluster'ı görselleştirmek için:

```bash
docker run -d \
  -p 9000:9000 \
  -e KAFKA_BROKERS=localhost:9092 \
  obsidiandynamics/kafdrop
```

Erişim: `http://localhost:9000`

## Local Development

### Docker Compose

```yaml
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
  
  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

### Topic Oluşturma

```bash
docker exec -it portfolio-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create \
  --topic portfolio.articles \
  --partitions 3 \
  --replication-factor 1
```

## Production Considerations

### High Availability

- **Replication Factor**: Minimum 3
- **Partitions**: Load'a göre ayarla
- **Brokers**: Minimum 3 broker

### Performance

- **Batch Size**: Producer batch size optimize et
- **Compression**: gzip veya snappy compression
- **Acks**: `all` (durability için)

### Security

- **SASL/SCRAM**: Authentication
- **TLS**: Encryption
- **ACLs**: Topic access control

## Best Practices

1. **Idempotent Consumers**: Consumer'lar idempotent olmalı
2. **Error Handling**: Retry ve DLQ kullan
3. **Monitoring**: Consumer lag'ı izle
4. **Partitioning**: Key-based partitioning kullan
5. **Schema Registry**: Event schema'ları için Avro/Protobuf kullan (opsiyonel)

