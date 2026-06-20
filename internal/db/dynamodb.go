package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Jhonni1000/ztna-proxy/internal/models"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamodbStore struct {
	dynamodbConn *dynamodb.Client
	myTableName  string
}

func InitialiseClient(ctx context.Context) *DynamodbStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Unable to load AWS default configuration, %w", err)
	}

	dynamodbClient := dynamodb.NewFromConfig(cfg)

	return &DynamodbStore{dynamodbClient, "ztna-policies"}
}

func (c *DynamodbStore) Create(ctx context.Context, policy *models.Policy) error {

	policy.Id, policy.CreatedAt, policy.UpdatedAt = uuid.New().String(), time.Now(), time.Now()
	dynamodbMap, err := attributevalue.MarshalMap(policy)
	if err != nil {
		return err
	}

	dynamoDBInput := dynamodb.PutItemInput{
		TableName: &c.myTableName,
		Item:      dynamodbMap,
	}

	_, err = c.dynamodbConn.PutItem(ctx, &dynamoDBInput)
	if err != nil {
		return err
	}

	return nil
}

func (c *DynamodbStore) Get(ctx context.Context, Id string) (*models.Policy, error) {
	lookupMap := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: Id},
	}

	dynamoDBInput := dynamodb.GetItemInput{
		TableName: &c.myTableName,
		Key:       lookupMap,
	}

	dynamodbOutput, err := c.dynamodbConn.GetItem(ctx, &dynamoDBInput)
	if err != nil {
		return nil, err
	}

	if len(dynamodbOutput.Item) == 0 || dynamodbOutput.Item == nil {
		return nil, fmt.Errorf("Policy not Found")
	}

	output := models.Policy{}
	err = attributevalue.UnmarshalMap(dynamodbOutput.Item, &output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (c *DynamodbStore) List(ctx context.Context) ([]*models.Policy, error) {

	queryInput := dynamodb.ScanInput{
		TableName: &c.myTableName,
	}

	dynamodbOutput, err := c.dynamodbConn.Scan(ctx, &queryInput)
	if err != nil {
		return nil, fmt.Errorf("Unable to query policies %w", err)
	}

	policies := []*models.Policy{}
	err = attributevalue.UnmarshalListOfMaps(dynamodbOutput.Items, &policies)
	if err != nil {
		return nil, fmt.Errorf("Error while Unmarshalling Map %w", err)
	}

	return policies, nil
}

func (c *DynamodbStore) Delete(ctx context.Context, Id string) error {
	lookupMap := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: Id},
	}

	dynamodbInput := dynamodb.DeleteItemInput{
		Key:       lookupMap,
		TableName: &c.myTableName,
	}

	_, err := c.dynamodbConn.DeleteItem(ctx, &dynamodbInput)
	if err != nil {
		return fmt.Errorf("Unable to delete input, %w", err)
	}
	return nil
}
