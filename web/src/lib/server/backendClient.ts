import type { DescService } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { env } from "$env/dynamic/private";

const transport = createConnectTransport({
  baseUrl: env.BACKEND_URL ?? "http://localhost:8080",
});

export function createBackendClient<T extends DescService>(service: T) {
  return createClient(service, transport);
}
