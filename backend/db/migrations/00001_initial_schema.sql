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

-- Deliberately no goose Down annotation/section below this line.
-- wrangler's "d1 migrations apply" (used to apply this same file to
-- production D1) has no rollback concept and just executes the whole
-- file verbatim, so a Down section here would DROP both tables
-- immediately after creating them in production. goose itself works
-- fine with a migration that has no Down section.
