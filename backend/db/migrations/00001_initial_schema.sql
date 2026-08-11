-- +goose Up
-- IF NOT EXISTS lets this migration double as a production baseline:
-- the existing prod DB already has these tables/indexes from Prisma
-- (byte-identical names/types, see backend/README.md), so running
-- `goose up` against it is a safe no-op that just records version 1 as
-- applied, instead of failing on "relation already exists".
CREATE TABLE IF NOT EXISTS "Follower" (
    "id" SERIAL PRIMARY KEY,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "publicKeyPem" TEXT NOT NULL,
    "following" BOOLEAN NOT NULL DEFAULT true,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "Follower_actorId_key" ON "Follower"("actorId");

CREATE TABLE IF NOT EXISTS "RelayConnection" (
    "id" SERIAL PRIMARY KEY,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "connected" BOOLEAN NOT NULL DEFAULT false,
    "lastAcceptedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "RelayConnection_actorId_key" ON "RelayConnection"("actorId");

-- +goose Down
DROP TABLE "RelayConnection";
DROP TABLE "Follower";
