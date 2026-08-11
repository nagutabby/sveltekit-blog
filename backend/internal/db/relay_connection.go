package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (q *Queries) GetRelayConnectionByActorID(ctx context.Context, actorID string) (RelayConnection, error) {
	out, err := q.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &q.relayConnectionTable,
		Key:       actorIDKey(actorID),
	})
	if err != nil {
		return RelayConnection{}, err
	}
	if out.Item == nil {
		return RelayConnection{}, ErrNotFound
	}
	var r RelayConnection
	if err := attributevalue.UnmarshalMap(out.Item, &r); err != nil {
		return RelayConnection{}, err
	}
	return r, nil
}

// ListRelayConnections returns every relay connection. Unlike the previous
// "ORDER BY id" Postgres query, Scan does not guarantee an order; no caller
// depends on ordering (broadcasts fan out to all relays regardless).
func (q *Queries) ListRelayConnections(ctx context.Context) ([]RelayConnection, error) {
	var items []RelayConnection
	var startKey map[string]types.AttributeValue
	for {
		out, err := q.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         &q.relayConnectionTable,
			ExclusiveStartKey: startKey,
		})
		if err != nil {
			return nil, err
		}
		var page []RelayConnection
		if err := attributevalue.UnmarshalListOfMaps(out.Items, &page); err != nil {
			return nil, err
		}
		items = append(items, page...)
		if out.LastEvaluatedKey == nil {
			return items, nil
		}
		startKey = out.LastEvaluatedKey
	}
}

type UpsertRelayConnectionAcceptedParams struct {
	ActorId string
	Inbox   string
}

// UpsertRelayConnectionAccepted implements the Accept(Follow/Subscribe)
// activity's persistence: create the row if it doesn't exist, or refresh
// inbox/lastAcceptedAt if it does.
func (q *Queries) UpsertRelayConnectionAccepted(ctx context.Context, arg UpsertRelayConnectionAcceptedParams) (RelayConnection, error) {
	now := nowTimestamp()
	out, err := q.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &q.relayConnectionTable,
		Key:       actorIDKey(arg.ActorId),
		UpdateExpression: strPtr(
			"SET Inbox = :inbox, Connected = :true, LastAcceptedAt = :now, UpdatedAt = :now, CreatedAt = if_not_exists(CreatedAt, :now)",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inbox": &types.AttributeValueMemberS{Value: arg.Inbox},
			":true":  &types.AttributeValueMemberBOOL{Value: true},
			":now":   &types.AttributeValueMemberS{Value: now},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return RelayConnection{}, err
	}
	var r RelayConnection
	if err := attributevalue.UnmarshalMap(out.Attributes, &r); err != nil {
		return RelayConnection{}, err
	}
	return r, nil
}
