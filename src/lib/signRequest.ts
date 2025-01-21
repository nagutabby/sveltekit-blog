import crypto from 'crypto';
import canonicalize from 'canonicalize';

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

export async function signCreateActivity(activity: any, privateKeyPem: string): Promise<object> {
  // 署名時刻の生成
  const created = new Date().toISOString();

  // 署名対象データの準備
  const dataToSign = {
    '@context': activity['@context'],
    type: activity.type,
    actor: activity.actor,
    object: activity.object,
    created: created
  };

  // データの正規化
  const canonicalizedData = canonicalize(dataToSign) as string;

  // 署名の生成
  const signature = crypto.sign(
    'sha256',
    Buffer.from(canonicalizedData),
    privateKeyPem
  );

  // LDSignatures形式での署名オブジェクト
  return {
    type: "RsaSignature2017",
    creator: `${activity.actor}#main-key`,
    created: created,
    signatureValue: signature.toString('base64')
  };
}
