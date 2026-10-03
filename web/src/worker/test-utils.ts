import type { WorkerEnv } from './bindings';
import type { Follower, RelayConnection } from './db';

export const siteBaseURL = 'https://blog.nagutabby.uk';

export class FakeD1 implements D1Database {
  readonly followers = new Map<string, Follower>();
  readonly relayConnections = new Map<string, RelayConnection>();
  failQueries = false;
  private nextFollowerID = 1;
  private nextRelayID = 1;

  prepare(query: string): FakeD1Statement {
    return new FakeD1Statement(this, query);
  }

  async batch<T = unknown>(statements: D1PreparedStatement[]): Promise<D1Result<T>[]> {
    return Promise.all(statements.map((statement) => statement.run<T>()));
  }

  async exec(_query: string): Promise<D1ExecResult> {
    if (this.failQueries) throw new Error('simulated D1 failure');
    return { count: 0, duration: 0 };
  }

  withSession(_constraintOrBookmark?: D1SessionBookmark | D1SessionConstraint): D1DatabaseSession {
    return {
      prepare: (query) => this.prepare(query),
      batch: <T = unknown>(statements: D1PreparedStatement[]) => this.batch<T>(statements),
      getBookmark: () => null
    };
  }

  async dump(): Promise<ArrayBuffer> {
    if (this.failQueries) throw new Error('simulated D1 failure');
    return new ArrayBuffer(0);
  }

  statement(query: string, parameters: unknown[]): unknown {
    if (this.failQueries) throw new Error('simulated D1 failure');
    if (query.includes('SELECT count(*) AS count FROM "Follower"')) {
      return { count: [...this.followers.values()].filter((row) => row.following).length };
    }
    if (query.includes('SELECT * FROM "Follower" WHERE "actorId" = ?')) {
      return this.followers.get(String(parameters[0])) ?? null;
    }
    if (query.includes('SELECT "actorId" FROM "Follower"')) {
      const offset = Number(parameters[1]);
      const limit = Number(parameters[0]);
      return [...this.followers.values()]
        .filter((row) => row.following)
        .sort((left, right) => left.id - right.id)
        .slice(offset, offset + limit)
        .map(({ actorId }) => ({ actorId }));
    }
    if (query.includes('SELECT * FROM "RelayConnection"')) {
      return [...this.relayConnections.values()].sort((left, right) => left.id - right.id);
    }
    if (query.includes('INSERT INTO "Follower"')) {
      const actorId = String(parameters[0]);
      const inbox = String(parameters[1]);
      const publicKeyPem = String(parameters[2]);
      const now = String(parameters[3]);
      const existing = this.followers.get(actorId);
      this.followers.set(actorId, {
        id: existing?.id ?? this.nextFollowerID++,
        actorId,
        inbox,
        publicKeyPem,
        following: true,
        createdAt: existing?.createdAt ?? now,
        updatedAt: now
      });
      return { changes: 1 };
    }
    if (query.includes('UPDATE "Follower"')) {
      const inbox = String(parameters[0]);
      const publicKeyPem = String(parameters[1]);
      const now = String(parameters[2]);
      const actorId = String(parameters[3]);
      const existing = this.followers.get(actorId);
      if (!existing) return { changes: 0 };
      this.followers.set(actorId, { ...existing, inbox, publicKeyPem, following: false, updatedAt: now });
      return { changes: 1 };
    }
    if (query.includes('INSERT INTO "RelayConnection"')) {
      const actorId = String(parameters[0]);
      const inbox = String(parameters[1]);
      const now = String(parameters[2]);
      const existing = this.relayConnections.get(actorId);
      this.relayConnections.set(actorId, {
        id: existing?.id ?? this.nextRelayID++,
        actorId,
        inbox,
        connected: true,
        lastAcceptedAt: now,
        createdAt: existing?.createdAt ?? now,
        updatedAt: now
      });
      return { changes: 1 };
    }
    throw new Error(`Unexpected SQL in Worker test: ${query}`);
  }
}

class FakeD1Statement implements D1PreparedStatement {
  private parameters: unknown[] = [];

  constructor(private readonly db: FakeD1, private readonly query: string) {}

  bind(...parameters: unknown[]): this {
    this.parameters = parameters;
    return this;
  }

  async first<T>(columnName?: string): Promise<T | null> {
    void columnName;
    return (this.db.statement(this.query, this.parameters) as T | null) ?? null;
  }

  async all<T>(): Promise<D1Result<T>> {
    const results = this.db.statement(this.query, this.parameters) as T[];
    return { results, success: true, meta: {} } as D1Result<T>;
  }

  async run<T>(): Promise<D1Result<T>> {
    const changes = this.db.statement(this.query, this.parameters) as { changes: number };
    return {
      results: [],
      success: true,
      meta: {
        duration: 0,
        size_after: 0,
        rows_read: 0,
        rows_written: changes.changes,
        last_row_id: 0,
        changed_db: changes.changes > 0,
        changes: changes.changes
      }
    } as D1Result<T>;
  }

  async raw<T = unknown[]>(options?: { columnNames?: false }): Promise<T[]>;
  async raw<T = unknown[]>(options: { columnNames: true }): Promise<[string[], ...T[]]>;
  async raw<T = unknown[]>(options?: { columnNames?: boolean }): Promise<T[] | [string[], ...T[]]> {
    const results = this.db.statement(this.query, this.parameters) as T[];
    if (options?.columnNames) return [[], ...results];
    return results;
  }
}

export function makeWorkerEnv(db = new FakeD1()): WorkerEnv {
  return {
    DB: db,
    ASSETS: {
      fetch: async () => new Response(null, { status: 404 }),
      connect: () => { throw new Error('Assets binding does not support sockets'); }
    },
    SITE_BASE_URL: siteBaseURL,
    ACTOR_PUBLIC_KEY_PEM: '',
    ACTOR_PRIVATE_KEY_PEM: '',
    FEDERATION_ADMIN_TOKEN: 'test-admin-token',
    EMAIL_API_TOKEN: 'test-email-token',
    FROM_ADDRESS: 'from@example.com',
    BCC_ADDRESS: 'bcc@example.com'
  };
}

export function seedFollower(db: FakeD1, actorId: string, publicKeyPem: string): Follower {
  const follower: Follower = {
    id: db.followers.size + 1,
    actorId,
    inbox: `${new URL(actorId).origin}/inbox`,
    publicKeyPem,
    following: true,
    createdAt: '2026-01-01T00:00:00.000Z',
    updatedAt: '2026-01-01T00:00:00.000Z'
  };
  db.followers.set(actorId, follower);
  return follower;
}
