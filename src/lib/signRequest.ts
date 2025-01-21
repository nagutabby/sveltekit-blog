import crypto from 'crypto';

export async function signRequest(url: string, method: string, body: string, privateKeyPem: string): Promise<Record<string, string>> {
  const parsedUrl = new URL(url);
  const host = parsedUrl.hostname; // hostnameを使用
  const date = new Date().toUTCString();

  // ダイジェストヘッダーの形式を修正
  const digest = `SHA-256=${crypto.createHash('sha256').update(body).digest('base64')}`;

  // シグネチャ文字列の構築を修正
  const signString = `(request-target): ${method.toLowerCase()} ${parsedUrl.pathname}\n` +
    `host: ${host}\n` +
    `date: ${date}\n` +
    `digest: ${digest}`;

  const signature = crypto.sign('sha256', Buffer.from(signString), privateKeyPem);

  const signatureHeader = [
    'keyId="https://blog.nagutabby.uk/actor#main-key"',
    'algorithm="rsa-sha256"',
    'headers="(request-target) host date digest"',
    `signature="${signature.toString('base64')}"`
  ].join(',');

  return {
    'Host': host,
    'Date': date,
    'Digest': digest,
    'Signature': signatureHeader,
    'Content-Type': 'application/activity+json',
    'Accept': 'application/activity+json'
  };
}
