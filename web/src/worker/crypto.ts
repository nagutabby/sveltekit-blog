import canonicalize from 'canonicalize';

const encoder = new TextEncoder();
const RSA_ALGORITHM: RsaHashedImportParams = { name: 'RSASSA-PKCS1-v1_5', hash: 'SHA-256' };
const RSA_SIGN_ALGORITHM: AlgorithmIdentifier = 'RSASSA-PKCS1-v1_5';
const RSA_ENCRYPTION_ALGORITHM = new Uint8Array([
  0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01, 0x05, 0x00
]);

export interface SignedRequestHeaders {
  date: string;
  digest: string;
  signature: string;
}

export interface LdSignature {
  type: 'RsaSignature2017';
  creator: string;
  created: string;
  signatureValue: string;
}

export function normalizePEM(value: string): string {
  return value.replaceAll('\\n', '\n');
}

function der(tag: number, bytes: Uint8Array): Uint8Array {
  const length: number[] = [];
  if (bytes.length < 0x80) {
    length.push(bytes.length);
  } else {
    let remaining = bytes.length;
    const encoded: number[] = [];
    while (remaining > 0) {
      encoded.unshift(remaining & 0xff);
      remaining >>>= 8;
    }
    length.push(0x80 | encoded.length, ...encoded);
  }
  return new Uint8Array([tag, ...length, ...bytes]);
}

function concatBytes(...parts: Uint8Array[]): Uint8Array {
  const result = new Uint8Array(parts.reduce((total, part) => total + part.length, 0));
  let offset = 0;
  for (const part of parts) {
    result.set(part, offset);
    offset += part.length;
  }
  return result;
}

function pemBlock(value: string): { label: string; der: Uint8Array } {
  const match = /-----BEGIN ([A-Z0-9 ]+)-----([\s\S]*?)-----END \1-----/.exec(normalizePEM(value));
  if (!match) throw new Error('federation: failed to decode PEM block');
  const label = match[1];
  const contents = match[2];
  if (!label || contents === undefined) throw new Error('federation: failed to decode PEM block');
  const encoded = contents.replace(/\s/g, '');
  const decoded = atob(encoded);
  return {
    label,
    der: Uint8Array.from(decoded, (character) => character.charCodeAt(0))
  };
}

function toArrayBuffer(bytes: Uint8Array): ArrayBuffer {
  const copy = new Uint8Array(bytes.length);
  copy.set(bytes);
  return copy.buffer;
}

function toBase64(bytes: Uint8Array): string {
  let binary = '';
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary);
}

function fromBase64(value: string): Uint8Array {
  const binary = atob(value);
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

async function importPrivateKey(pem: string): Promise<CryptoKey> {
  const block = pemBlock(pem);
  let keyData: Uint8Array;
  if (block.label === 'PRIVATE KEY') {
    keyData = block.der;
  } else if (block.label === 'RSA PRIVATE KEY') {
    const version = new Uint8Array([0x02, 0x01, 0x00]);
    keyData = der(0x30, concatBytes(version, RSA_ENCRYPTION_ALGORITHM, der(0x04, block.der)));
  } else {
    throw new Error('federation: private key is not PKCS#1 or PKCS#8 RSA');
  }
  return crypto.subtle.importKey('pkcs8', toArrayBuffer(keyData), RSA_ALGORITHM, false, ['sign']);
}

async function importPublicKey(pem: string): Promise<CryptoKey> {
  const block = pemBlock(pem);
  let keyData: Uint8Array;
  if (block.label === 'PUBLIC KEY') {
    keyData = block.der;
  } else if (block.label === 'RSA PUBLIC KEY') {
    keyData = der(0x30, concatBytes(RSA_ENCRYPTION_ALGORITHM, der(0x03, concatBytes(new Uint8Array([0]), block.der))));
  } else {
    throw new Error('federation: public key is not PKIX or PKCS#1 RSA');
  }
  return crypto.subtle.importKey('spki', toArrayBuffer(keyData), RSA_ALGORITHM, false, ['verify']);
}

export async function sha256Base64(value: string): Promise<string> {
  const digest = await crypto.subtle.digest('SHA-256', encoder.encode(value));
  return toBase64(new Uint8Array(digest));
}

export async function signHTTPRequest(
  targetURL: string,
  method: string,
  body: string,
  keyID: string,
  privateKeyPEM: string
): Promise<SignedRequestHeaders> {
  const parsed = new URL(targetURL);
  const date = new Date().toUTCString();
  const digest = `SHA-256=${await sha256Base64(body)}`;
  const requestTarget = `${parsed.pathname}${parsed.search}`;
  const signingString = [
    `(request-target): ${method.toLowerCase()} ${requestTarget}`,
    `host: ${parsed.host}`,
    `date: ${date}`,
    `digest: ${digest}`
  ].join('\n');
  const key = await importPrivateKey(privateKeyPEM);
  const signature = await crypto.subtle.sign(RSA_SIGN_ALGORITHM, key, encoder.encode(signingString));
  return {
    date,
    digest,
    signature: `keyId="${keyID}",algorithm="rsa-sha256",headers="(request-target) host date digest",signature="${toBase64(new Uint8Array(signature))}"`
  };
}

function parseSignatureHeader(header: string): Map<string, string> {
  const params = new Map<string, string>();
  for (const match of header.matchAll(/(\w+)="([^"]*)"/g)) {
    const key = match[1];
    const value = match[2];
    if (key !== undefined && value !== undefined) params.set(key, value);
  }
  if (!params.size || !params.get('keyId') || !params.get('signature')) {
    throw new Error('federation: malformed Signature header');
  }
  return params;
}

function signingStringForRequest(request: Request, headerList: string[]): string {
  const url = new URL(request.url);
  const requestTarget = `${url.pathname}${url.search}`;
  const host = request.headers.get('host') ?? url.host;
  return headerList.map((header) => {
    const lower = header.toLowerCase();
    let value: string | null;
    if (lower === '(request-target)') value = `${request.method.toLowerCase()} ${requestTarget}`;
    else if (lower === 'host') value = host;
    else value = request.headers.get(header);
    if (value === null || value === '') throw new Error(`federation: signed header "${header}" is missing from the request`);
    return `${lower}: ${value}`;
  }).join('\n');
}

export async function verifyHTTPRequest(
  request: Request,
  body: Uint8Array,
  publicKeyPEM: string
): Promise<void> {
  const signatureHeader = request.headers.get('Signature');
  if (!signatureHeader) throw new Error('federation: missing Signature header');

  const params = parseSignatureHeader(signatureHeader);
  const algorithm = params.get('algorithm');
  if (algorithm && algorithm !== 'rsa-sha256') throw new Error(`federation: unsupported signature algorithm "${algorithm}"`);

  const headerList = (params.get('headers') ?? '').split(/\s+/).filter(Boolean);
  if (!headerList.length) headerList.push('date');

  const dateHeader = request.headers.get('Date');
  if (!dateHeader) throw new Error('federation: missing Date header');
  const parsedDate = Date.parse(dateHeader);
  if (!Number.isFinite(parsedDate)) throw new Error('federation: invalid Date header');
  if (Math.abs(Date.now() - parsedDate) > 60 * 60 * 1000) {
    throw new Error('federation: Date header outside acceptable window');
  }

  if (body.length > 0) {
    if (!headerList.some((header) => header.toLowerCase() === 'digest')) {
      throw new Error('federation: Digest header must be part of the signed headers');
    }
    const digest = request.headers.get('Digest');
    if (!digest) throw new Error('federation: missing Digest header');
    const bodyDigest = await crypto.subtle.digest('SHA-256', toArrayBuffer(body));
    const expected = `SHA-256=${toBase64(new Uint8Array(bodyDigest))}`;
    if (digest !== expected) throw new Error('federation: Digest header does not match body');
  }

  const key = await importPublicKey(publicKeyPEM);
  const signingString = signingStringForRequest(request, headerList);
  const verified = await crypto.subtle.verify(
    RSA_SIGN_ALGORITHM,
    key,
    toArrayBuffer(fromBase64(params.get('signature') ?? '')),
    encoder.encode(signingString)
  );
  if (!verified) throw new Error('federation: signature verification failed');
}

export async function signActivity(
  activity: {
    '@context': string[];
    type: string;
    actor: string;
    object: Record<string, unknown>;
  },
  privateKeyPEM: string
): Promise<LdSignature> {
  const created = new Date().toISOString().replace(/Z$/, '').slice(0, 23) + 'Z';
  const unsigned = {
    '@context': activity['@context'],
    type: activity.type,
    actor: activity.actor,
    object: activity.object,
    created
  };
  const canonical = canonicalize(unsigned);
  if (canonical === undefined) throw new Error('federationadmin: canonicalizing activity failed');

  const key = await importPrivateKey(privateKeyPEM);
  const signature = await crypto.subtle.sign(RSA_SIGN_ALGORITHM, key, encoder.encode(canonical));
  return {
    type: 'RsaSignature2017',
    creator: `${activity.actor}#main-key`,
    created,
    signatureValue: toBase64(new Uint8Array(signature))
  };
}

export async function constantTimeTokenMatches(provided: string, expected: string): Promise<boolean> {
  const [providedHash, expectedHash] = await Promise.all([
    crypto.subtle.digest('SHA-256', encoder.encode(provided)),
    crypto.subtle.digest('SHA-256', encoder.encode(expected))
  ]);
  const subtle = crypto.subtle as SubtleCrypto & { timingSafeEqual?: (left: ArrayBuffer, right: ArrayBuffer) => boolean };
  if (subtle.timingSafeEqual) return subtle.timingSafeEqual(providedHash, expectedHash);
  const left = new Uint8Array(providedHash);
  const right = new Uint8Array(expectedHash);
  let difference = 0;
  for (let index = 0; index < left.length; index += 1) {
    difference |= (left[index] ?? 0) ^ (right[index] ?? 0);
  }
  return difference === 0;
}
