-- +goose Up
CREATE TABLE "Follower" (
    "id" SERIAL PRIMARY KEY,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "publicKeyPem" TEXT NOT NULL,
    "following" BOOLEAN NOT NULL DEFAULT true,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL
);

CREATE UNIQUE INDEX "Follower_actorId_key" ON "Follower"("actorId");

CREATE TABLE "RelayConnection" (
    "id" SERIAL PRIMARY KEY,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "connected" BOOLEAN NOT NULL DEFAULT false,
    "lastAcceptedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL
);

CREATE UNIQUE INDEX "RelayConnection_actorId_key" ON "RelayConnection"("actorId");

-- +goose Down
DROP TABLE "RelayConnection";
DROP TABLE "Follower";
