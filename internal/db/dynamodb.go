package db

import (
	"context"
	"log"
	"time"

	"github.com/Jhonni1000/ztna-proxy/internal/models"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type DynamodbStore struct {
	dynamodbConn *dynamodb.Client
}

func InitialiseClient(ctx context.Context) DynamodbStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load AWS default configuration, %s", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	return DynamodbStore{dynamodbClient}
}

func (c *DynamodbStore) Create(ctx context.Context, policy *models.Policy) error {

	policy.Id, policy.CreatedAt, policy.UpdatedAt = uuid.New().String(), time.Now(), time.Now()
	dynamodbMap, err := attributevalue.MarshalMap(policy)
	if err != nil {
		return err
	}

	var myTableName = "ztna-policies"
	dynamoDBInput := dynamodb.PutItemInput{
		TableName: &myTableName,
		Item:      dynamodbMap,
	}

	_, err = c.dynamodbConn.PutItem(ctx, &dynamoDBInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *DynamodbStore) Get(ctx context.Context, id string) (*models.Policy, error) {

}

func (c *DynamodbStore) List(ctx context.Context) ([]*models.Policy, error) {

}

func (c *DynamodbStore) Delete(ctx context.Context, id string) error {

}
