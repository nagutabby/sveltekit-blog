package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (q *Queries) CountActiveFollowers(ctx context.Context) (int64, error) {
	var count int64
	var startKey map[string]types.AttributeValue
	for {
		out, err := q.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        &q.followerTable,
			FilterExpression: strPtr("Following = :true"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":true": &types.AttributeValueMemberBOOL{Value: true},
			},
			Select:            types.SelectCount,
			ExclusiveStartKey: startKey,
		})
		if err != nil {
			return 0, err
		}
		count += int64(out.Count)
		if out.LastEvaluatedKey == nil {
			return count, nil
		}
		startKey = out.LastEvaluatedKey
	}
}

func (q *Queries) GetFollowerByActorID(ctx context.Context, actorID string) (Follower, error) {
	out, err := q.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &q.followerTable,
		Key:       actorIDKey(actorID),
	})
	if err != nil {
		return Follower{}, err
	}
	if out.Item == nil {
		return Follower{}, ErrNotFound
	}
	var f Follower
	if err := attributevalue.UnmarshalMap(out.Item, &f); err != nil {
		return Follower{}, err
	}
	return f, nil
}

type UpsertFollowerParams struct {
	ActorId      string
	Inbox        string
	PublicKeyPem string
}

// UpsertFollower implements the Follow activity's persistence: create the
// row if it doesn't exist, or refresh inbox/publicKeyPem and mark it
// following again if it does.
func (q *Queries) UpsertFollower(ctx context.Context, arg UpsertFollowerParams) (Follower, error) {
	now := nowTimestamp()
	out, err := q.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &q.followerTable,
		Key:       actorIDKey(arg.ActorId),
		UpdateExpression: strPtr(
			"SET Inbox = :inbox, PublicKeyPem = :publicKeyPem, Following = :true, UpdatedAt = :now, CreatedAt = if_not_exists(CreatedAt, :now)",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inbox":        &types.AttributeValueMemberS{Value: arg.Inbox},
			":publicKeyPem": &types.AttributeValueMemberS{Value: arg.PublicKeyPem},
			":true":         &types.AttributeValueMemberBOOL{Value: true},
			":now":          &types.AttributeValueMemberS{Value: now},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return Follower{}, err
	}
	var f Follower
	if err := attributevalue.UnmarshalMap(out.Attributes, &f); err != nil {
		return Follower{}, err
	}
	return f, nil
}

type UnfollowByActorIDParams struct {
	ActorId      string
	Inbox        string
	PublicKeyPem string
}

// UnfollowByActorID implements the Undo(Follow) activity's persistence.
// It fails if the actor was never a follower, matching the previous
// Postgres UPDATE ... RETURNING behavior (0 rows affected -> ErrNoRows).
func (q *Queries) UnfollowByActorID(ctx context.Context, arg UnfollowByActorIDParams) (Follower, error) {
	now := nowTimestamp()
	out, err := q.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:           &q.followerTable,
		Key:                 actorIDKey(arg.ActorId),
		ConditionExpression: strPtr("attribute_exists(ActorId)"),
		UpdateExpression: strPtr(
			"SET Inbox = :inbox, PublicKeyPem = :publicKeyPem, Following = :false, UpdatedAt = :now",
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inbox":        &types.AttributeValueMemberS{Value: arg.Inbox},
			":publicKeyPem": &types.AttributeValueMemberS{Value: arg.PublicKeyPem},
			":false":        &types.AttributeValueMemberBOOL{Value: false},
			":now":          &types.AttributeValueMemberS{Value: now},
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return Follower{}, err
	}
	var f Follower
	if err := attributevalue.UnmarshalMap(out.Attributes, &f); err != nil {
		return Follower{}, err
	}
	return f, nil
}
