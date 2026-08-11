import type { DescService } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

const transport = createConnectTransport({
  baseUrl: "/"
});

export function createBackendClient<T extends DescService>(service: T) {
  return createClient(service, transport);
}
