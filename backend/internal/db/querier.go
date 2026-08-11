package db

import (
	"context"
)

type Querier interface {
	CountActiveFollowers(ctx context.Context) (int64, error)
	GetFollowerByActorID(ctx context.Context, actorID string) (Follower, error)
	GetRelayConnectionByActorID(ctx context.Context, actorID string) (RelayConnection, error)
	ListRelayConnections(ctx context.Context) ([]RelayConnection, error)
	// Undo(Follow)と同じ挙動: 既存のフォロワーをfollowing=falseにする。レコードが存在しない場合はエラーになる。
	UnfollowByActorID(ctx context.Context, arg UnfollowByActorIDParams) (Follower, error)
	// Followと同じ挙動: 既存レコードがあればfollowing/inbox/publicKeyPemを更新し、なければ作成する。
	UpsertFollower(ctx context.Context, arg UpsertFollowerParams) (Follower, error)
	// Accept(Follow/Subscribe)と同じ挙動: リレーからのAcceptを記録する。
	UpsertRelayConnectionAccepted(ctx context.Context, arg UpsertRelayConnectionAcceptedParams) (RelayConnection, error)
}

var _ Querier = (*Queries)(nil)
