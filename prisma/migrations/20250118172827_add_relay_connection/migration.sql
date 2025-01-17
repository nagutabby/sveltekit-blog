-- CreateTable
CREATE TABLE "RelayConnection" (
    "id" SERIAL NOT NULL,
    "actorId" TEXT NOT NULL,
    "connected" BOOLEAN NOT NULL DEFAULT false,
    "lastAcceptedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "RelayConnection_pkey" PRIMARY KEY ("id")
);
