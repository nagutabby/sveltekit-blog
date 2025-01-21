import crypto from 'crypto';

export async function signRequest(url: string, method: string, body: string, privateKeyPem: string): Promise<Record<string, string>> {
  const host = new URL(url).host;
  const date = new Date().toUTCString();
  const digest = crypto.createHash('sha256').update(body).digest('base64');
  const digestHeader = `sha-256=${digest}`;

  const signString = [
    `(request-target): ${method.toLowerCase()} ${new URL(url).pathname}`,
    `host: ${host}`,
    `date: ${date}`,
    `digest: ${digestHeader}`
  ].join('\n');

  const signer = crypto.createSign('sha256');
  signer.update(signString);
  const signature = signer.sign(privateKeyPem, 'base64');

  const signatureHeader = [
    'keyId="https://blog.nagutabby.uk/actor#main-key"',
    'algorithm="rsa-sha256"',
    'headers="(request-target) host date digest"',
    `signature="${signature}"`
  ].join(',');

  return {
    'Host': host,
    'Date': date,
    'Digest': digestHeader,
    'Signature': signatureHeader,
    'Content-Type': 'application/activity+json'
  };
}
