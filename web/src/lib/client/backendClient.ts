import type { DescService } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

// SSG化前は同一オリジン(Cloudflare Workerがbackendへ転送)だったため相対パスで
// 十分だった。SSG化後はblog.nagutabby.uk(静的サイト)からapi.nagutabby.uk
// (Lambda)への別オリジン呼び出しになるため、ビルド時にPUBLIC_API_BASE_URLで
// 明示する。未設定時は相対パス(同一オリジン、ローカル開発時のvite proxy等)のまま。
const baseUrl = import.meta.env.PUBLIC_API_BASE_URL || "/";

const transport = createConnectTransport({
  baseUrl
});

export function createBackendClient<T extends DescService>(service: T) {
  return createClient(service, transport);
}
