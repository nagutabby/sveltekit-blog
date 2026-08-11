-- name: UpsertRelayConnectionAccepted :one
-- Accept(Follow/Subscribe)と同じ挙動: リレーからのAcceptを記録する。
INSERT INTO "RelayConnection" ("actorId", "inbox", "connected", "lastAcceptedAt", "updatedAt")
VALUES ($1, $2, true, now(), now())
ON CONFLICT ("actorId") DO UPDATE SET
    "connected" = true,
    "inbox" = EXCLUDED."inbox",
    "lastAcceptedAt" = now(),
    "updatedAt" = now()
RETURNING *;

-- name: ListRelayConnections :many
SELECT * FROM "RelayConnection" ORDER BY "id";

-- name: GetRelayConnectionByActorID :one
SELECT * FROM "RelayConnection" WHERE "actorId" = $1;
