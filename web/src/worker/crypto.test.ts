import { webcrypto } from 'node:crypto';
import canonicalize from 'canonicalize';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  constantTimeTokenMatches,
  normalizePEM,
  sha256Base64,
  signActivity,
  signHTTPRequest,
  verifyHTTPRequest
} from './crypto';

vi.stubGlobal('crypto', webcrypto);

function toPEM(label: string, bytes: ArrayBuffer): string {
  const binary = Array.from(new Uint8Array(bytes), (byte) => String.fromCharCode(byte)).join('');
  const encoded = btoa(binary).match(/.{1,64}/g)?.join('\n') ?? '';
  return `-----BEGIN ${label}-----\n${encoded}\n-----END ${label}-----`;
}

async function makeKeyPair() {
  const pair = await crypto.subtle.generateKey(
    { name: 'RSASSA-PKCS1-v1_5', modulusLength: 2048, publicExponent: new Uint8Array([1, 0, 1]), hash: 'SHA-256' },
    true,
    ['sign', 'verify']
  );
  const [privateBytes, publicBytes] = await Promise.all([
    crypto.subtle.exportKey('pkcs8', pair.privateKey),
    crypto.subtle.exportKey('spki', pair.publicKey)
  ]);
  return {
    privateKeyPEM: toPEM('PRIVATE KEY', privateBytes),
    publicKeyPEM: toPEM('PUBLIC KEY', publicBytes),
    publicKey: pair.publicKey
  };
}

describe('Worker cryptography', () => {
  beforeEach(() => {
    vi.stubGlobal('crypto', webcrypto);
  });

  it('normalizes escaped PEM line breaks and hashes request bodies', async () => {
    expect(normalizePEM('one\\ntwo')).toBe('one\ntwo');
    await expect(sha256Base64('hello')).resolves.toBe('LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=');
  });

  it('signs and verifies an HTTP request, including a query in request-target', async () => {
    const keys = await makeKeyPair();
    const target = 'https://relay.example/inbox?shared=true';
    const body = '{"type":"Follow"}';
    const headers = await signHTTPRequest(target, 'POST', body, 'https://blog.example/actor#main-key', keys.privateKeyPEM);
    const request = new Request(target, {
      method: 'POST',
      headers: { Date: headers.date, Digest: headers.digest, Signature: headers.signature },
      body
    });

    await expect(verifyHTTPRequest(request, new TextEncoder().encode(body), keys.publicKeyPEM)).resolves.toBeUndefined();
    await expect(verifyHTTPRequest(request, new TextEncoder().encode('{"type":"Delete"}'), keys.publicKeyPEM)).rejects.toThrow(/Digest/);
  });

  it('creates a verifiable JCS LD-Signature', async () => {
    const keys = await makeKeyPair();
    const activity = {
      '@context': ['https://www.w3.org/ns/activitystreams', 'https://w3id.org/security/v1'],
      type: 'Create',
      actor: 'https://blog.example/actor',
      object: { id: 'https://blog.example/api/articles/post', type: 'Note', name: 'title' }
    };
    const signature = await signActivity(activity, keys.privateKeyPEM);
    const unsigned = { ...activity, created: signature.created };
    const canonical = canonicalize(unsigned);
    expect(signature.type).toBe('RsaSignature2017');
    expect(canonical).toBeDefined();
    const signatureBytes = Uint8Array.from(atob(signature.signatureValue), (character) => character.charCodeAt(0));
    await expect(crypto.subtle.verify(
      'RSASSA-PKCS1-v1_5',
      keys.publicKey,
      signatureBytes,
      new TextEncoder().encode(canonical)
    )).resolves.toBe(true);
  });

  it('compares admin tokens and rejects a different token', async () => {
    await expect(constantTimeTokenMatches('secret', 'secret')).resolves.toBe(true);
    await expect(constantTimeTokenMatches('wrong', 'secret')).resolves.toBe(false);
  });
});
