package federation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// remoteActor is the subset of an ActivityPub Actor document this codebase
// needs to countersign and reply to activities.
type remoteActor struct {
	Inbox     string `json:"inbox"`
	PublicKey struct {
		ID           string `json:"id"`
		Owner        string `json:"owner"`
		PublicKeyPem string `json:"publicKeyPem"`
	} `json:"publicKey"`
}

var errActorFetchFailed = errors.New("federation: could not fetch actor information")

func fetchActor(ctx context.Context, client *http.Client, actorURL string) (*remoteActor, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, actorURL, nil)
	if err != nil {
		return nil, fmt.Errorf("federation: building actor request: %w", err)
	}
	req.Header.Set("Accept", "application/activity+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errActorFetchFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errActorFetchFailed
	}

	var actor remoteActor
	if err := json.NewDecoder(resp.Body).Decode(&actor); err != nil {
		return nil, errActorFetchFailed
	}

	if actor.Inbox == "" || actor.PublicKey.PublicKeyPem == "" {
		return nil, errActorFetchFailed
	}

	return &actor, nil
}
