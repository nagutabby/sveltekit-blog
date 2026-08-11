-- name: UpsertFollower :one
-- Followと同じ挙動: 既存レコードがあればfollowing/inbox/publicKeyPemを更新し、なければ作成する。
INSERT INTO "Follower" ("actorId", "inbox", "publicKeyPem", "following", "updatedAt")
VALUES ($1, $2, $3, true, now())
ON CONFLICT ("actorId") DO UPDATE SET
    "following" = true,
    "inbox" = EXCLUDED."inbox",
    "publicKeyPem" = EXCLUDED."publicKeyPem",
    "updatedAt" = now()
RETURNING *;

-- name: UnfollowByActorID :one
-- Undo(Follow)と同じ挙動: 既存のフォロワーをfollowing=falseにする。レコードが存在しない場合はエラーになる。
UPDATE "Follower"
SET "following" = false,
    "inbox" = $2,
    "publicKeyPem" = $3,
    "updatedAt" = now()
WHERE "actorId" = $1
RETURNING *;

-- name: CountActiveFollowers :one
SELECT count(*) FROM "Follower" WHERE "following" = true;

-- name: GetFollowerByActorID :one
SELECT * FROM "Follower" WHERE "actorId" = $1;
