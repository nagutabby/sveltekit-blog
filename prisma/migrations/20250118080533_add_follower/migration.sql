-- CreateTable
CREATE TABLE "Follower" (
    "id" SERIAL NOT NULL,
    "actorId" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "publicKeyPem" TEXT NOT NULL,
    "following" BOOLEAN NOT NULL DEFAULT true,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Follower_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Follower_actorId_key" ON "Follower"("actorId");
