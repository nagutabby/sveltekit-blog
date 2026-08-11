package db

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// FollowerTableDefinition and RelayConnectionTableDefinition are the single
// source of truth for the tables' key schema, shared by integration tests
// (against dynamodb-local) and the data migration CLI. The CDK DynamoDbStack
// defines the same key schema independently for the real AWS tables.
func FollowerTableDefinition(tableName string) *dynamodb.CreateTableInput {
	return actorIDTableDefinition(tableName)
}

func RelayConnectionTableDefinition(tableName string) *dynamodb.CreateTableInput {
	return actorIDTableDefinition(tableName)
}

func actorIDTableDefinition(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		TableName: &tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: strPtr("ActorId"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: strPtr("ActorId"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	}
}

func actorIDKey(actorID string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"ActorId": &types.AttributeValueMemberS{Value: actorID},
	}
}

func strPtr(s string) *string {
	return &s
}
