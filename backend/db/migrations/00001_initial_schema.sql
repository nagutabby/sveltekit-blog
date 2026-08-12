-- +goose Up
CREATE TABLE IF NOT EXISTS "Follower" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "publicKeyPem" TEXT NOT NULL,
    "following" BOOLEAN NOT NULL DEFAULT 1,
    "createdAt" TEXT NOT NULL,
    "updatedAt" TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "Follower_actorId_key" ON "Follower"("actorId");

CREATE TABLE IF NOT EXISTS "RelayConnection" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "connected" BOOLEAN NOT NULL DEFAULT 0,
    "lastAcceptedAt" TEXT,
    "createdAt" TEXT NOT NULL,
    "updatedAt" TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "RelayConnection_actorId_key" ON "RelayConnection"("actorId");

-- +goose Down
DROP TABLE "RelayConnection";
DROP TABLE "Follower";
