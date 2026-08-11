import { describe, expect, it } from "vitest";
import { createClient, createRouterTransport } from "@connectrpc/connect";
import { HealthService } from "./health_pb";

describe("HealthService generated stubs", () => {
  it("wires a client to a handler via the generated service descriptor", async () => {
    const transport = createRouterTransport(({ service }) => {
      service(HealthService, {
        check: () => ({ status: "ok" }),
      });
    });

    const client = createClient(HealthService, transport);
    const response = await client.check({});

    expect(response.status).toBe("ok");
  });
});
