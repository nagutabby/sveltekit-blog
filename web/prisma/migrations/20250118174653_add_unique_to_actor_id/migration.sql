/*
  Warnings:

  - A unique constraint covering the columns `[actorId]` on the table `RelayConnection` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "RelayConnection_actorId_key" ON "RelayConnection"("actorId");
