export function generateKeyPair() {
  const kp = nacl.box.keyPair();
  return {
    publicKey: kp.publicKey,
    privateKey: kp.secretKey
  };
}

export function deriveDH(priv, pub) {
  return nacl.scalarMult(priv, pub);
}

export async function deriveSharedSecret({ EK_A_priv, IK_B, SPK_B, OTPK_B }) {
  const dh1 = deriveDH(EK_A_priv, IK_B);
  const dh2 = deriveDH(EK_A_priv, SPK_B);
  const dh3 = deriveDH(EK_A_priv, OTPK_B);

  const combined = new Uint8Array(dh1.length * 3);
  combined.set(dh1, 0);
  combined.set(dh2, dh1.length);
  combined.set(dh3, dh1.length * 2);

  const hashBuffer = await crypto.subtle.digest("SHA-256", combined);
  return new Uint8Array(hashBuffer);
}

export function aesEncrypt(plaintext, key) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const encoder = new TextEncoder();
  const encoded = encoder.encode(plaintext);

  return crypto.subtle.importKey(
    "raw", key.slice(0, 32), { name: "AES-GCM" }, false, ["encrypt"]
  ).then(importedKey => {
    return crypto.subtle.encrypt({ name: "AES-GCM", iv }, importedKey, encoded)
      .then(cipher => ({
        ciphertext: new Uint8Array(cipher),
        iv
      }));
  });
}

export function base64ToBytes(str) {
  return Uint8Array.from(atob(str), c => c.charCodeAt(0));
}
export function bytesToBase64(bytes) {
  return btoa(String.fromCharCode(...bytes));
}
