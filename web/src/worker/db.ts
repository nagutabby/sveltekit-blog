export interface Follower {
  id: number;
  actorId: string;
  inbox: string;
  publicKeyPem: string;
  following: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface RelayConnection {
  id: number;
  actorId: string;
  inbox: string;
  connected: boolean;
  lastAcceptedAt: string | null;
  createdAt: string;
  updatedAt: string;
}

export async function countActiveFollowers(db: D1Database): Promise<number> {
  const row = await db
    .prepare('SELECT count(*) AS count FROM "Follower" WHERE "following" = 1')
    .first<{ count: number }>();
  return row?.count ?? 0;
}

export async function listActiveFollowerActorIDs(
  db: D1Database,
  limit: number,
  offset: number
): Promise<string[]> {
  const result = await db
    .prepare(
      'SELECT "actorId" FROM "Follower" WHERE "following" = 1 ORDER BY "id" LIMIT ? OFFSET ?'
    )
    .bind(limit, offset)
    .all<{ actorId: string }>();
  return result.results.map(({ actorId }) => actorId);
}

export async function getFollowerByActorID(
  db: D1Database,
  actorID: string
): Promise<Follower | null> {
  return db
    .prepare('SELECT * FROM "Follower" WHERE "actorId" = ?')
    .bind(actorID)
    .first<Follower>();
}

export async function upsertFollower(
  db: D1Database,
  params: Pick<Follower, 'actorId' | 'inbox' | 'publicKeyPem'> & { now: string }
): Promise<void> {
  await db
    .prepare(
      `INSERT INTO "Follower" ("actorId", "inbox", "publicKeyPem", "following", "createdAt", "updatedAt")
       VALUES (?, ?, ?, 1, ?, ?)
       ON CONFLICT ("actorId") DO UPDATE SET
         "following" = 1,
         "inbox" = excluded."inbox",
         "publicKeyPem" = excluded."publicKeyPem",
         "updatedAt" = excluded."updatedAt"`
    )
    .bind(params.actorId, params.inbox, params.publicKeyPem, params.now, params.now)
    .run();
}

export async function unfollowByActorID(
  db: D1Database,
  params: Pick<Follower, 'actorId' | 'inbox' | 'publicKeyPem'> & { now: string }
): Promise<void> {
  const result = await db
    .prepare(
      `UPDATE "Follower"
       SET "following" = 0, "inbox" = ?, "publicKeyPem" = ?, "updatedAt" = ?
       WHERE "actorId" = ?`
    )
    .bind(params.inbox, params.publicKeyPem, params.now, params.actorId)
    .run();
  if ((result.meta.changes ?? 0) === 0) throw new Error('follower not found');
}

export async function listRelayConnections(db: D1Database): Promise<RelayConnection[]> {
  const result = await db
    .prepare('SELECT * FROM "RelayConnection" ORDER BY "id"')
    .all<RelayConnection>();
  return result.results;
}

export async function upsertRelayConnectionAccepted(
  db: D1Database,
  params: { actorId: string; inbox: string; now: string }
): Promise<void> {
  await db
    .prepare(
      `INSERT INTO "RelayConnection" ("actorId", "inbox", "connected", "lastAcceptedAt", "createdAt", "updatedAt")
       VALUES (?, ?, 1, ?, ?, ?)
       ON CONFLICT ("actorId") DO UPDATE SET
         "connected" = 1,
         "inbox" = excluded."inbox",
         "lastAcceptedAt" = excluded."lastAcceptedAt",
         "updatedAt" = excluded."updatedAt"`
    )
    .bind(params.actorId, params.inbox, params.now, params.now, params.now)
    .run();
}
