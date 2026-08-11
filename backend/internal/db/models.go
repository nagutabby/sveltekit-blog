package db

// Follower is a remote ActivityPub actor following this blog's actor.
// CreatedAt/UpdatedAt are RFC3339Nano timestamps (sortable as plain strings).
type Follower struct {
	ActorId      string `dynamodbav:"ActorId"`
	Inbox        string `dynamodbav:"Inbox"`
	PublicKeyPem string `dynamodbav:"PublicKeyPem"`
	Following    bool   `dynamodbav:"Following"`
	CreatedAt    string `dynamodbav:"CreatedAt"`
	UpdatedAt    string `dynamodbav:"UpdatedAt"`
}

// RelayConnection is a relay this blog's actor has subscribed to.
// LastAcceptedAt is empty until the relay has sent an Accept.
type RelayConnection struct {
	ActorId        string `dynamodbav:"ActorId"`
	Inbox          string `dynamodbav:"Inbox"`
	Connected      bool   `dynamodbav:"Connected"`
	LastAcceptedAt string `dynamodbav:"LastAcceptedAt"`
	CreatedAt      string `dynamodbav:"CreatedAt"`
	UpdatedAt      string `dynamodbav:"UpdatedAt"`
}
