import { useEffect, useState } from "react";
import nacl from "tweetnacl";

function base64ToBytes(str: string): Uint8Array {
  return Uint8Array.from(atob(str), (c) => c.charCodeAt(0));
}

function bytesToBase64(bytes: Uint8Array): string {
  return btoa(String.fromCharCode(...bytes));
}

type KeyPair = {
  private: string;
  public: string;
};

type RecipientKeys = {
    identity_key: KeyPair;
    signed_pre_key: KeyPair;
    one_time_keys: KeyPair[];
  
    identity_key_signer: KeyPair;
    spk_signature: string;
};

export default function X3DHReceiver() {
  const [recipient, setRecipient] = useState<RecipientKeys | null>(null);
  const [ekAPub, setEkAPub] = useState<string>("");
  const [iv, setIv] = useState<string>("");
  const [ciphertext, setCiphertext] = useState<string>("");
  const [decrypted, setDecrypted] = useState<string>("");

  useEffect(() => {
    fetch("http://localhost:8002/keys")
      .then((res) => res.json())
      .then(setRecipient);
  }, []);

  const decrypt = async () => {
    if (!recipient) return;

    const EK_A = base64ToBytes(ekAPub);
    const IV = base64ToBytes(iv);
    const CT = base64ToBytes(ciphertext);

    const IK_B_priv = base64ToBytes(recipient.identity_key.private);
    const SPK_B_priv = base64ToBytes(recipient.signed_pre_key.private);
    const OTPK_B_priv = base64ToBytes(recipient.one_time_keys[0].private);

    const dh = (priv: Uint8Array, pub: Uint8Array) => nacl.scalarMult(priv, pub);
    const dh1 = dh(IK_B_priv, EK_A);
    const dh2 = dh(SPK_B_priv, EK_A);
    const dh3 = dh(OTPK_B_priv, EK_A);

    const combined = new Uint8Array([...dh1, ...dh2, ...dh3]);
    const hashBuffer = await crypto.subtle.digest("SHA-256", combined);
    const sharedSecret = new Uint8Array(hashBuffer);

    try {
      const key = await crypto.subtle.importKey("raw", sharedSecret, { name: "AES-GCM" }, false, ["decrypt"]);
      const decryptedBuffer = await crypto.subtle.decrypt({ name: "AES-GCM", iv: IV }, key, CT);
      const plaintext = new TextDecoder().decode(decryptedBuffer);
      setDecrypted(plaintext);
    } catch (e) {
      setDecrypted("❌ Ошибка при расшифровке: " + (e as Error).message);
    }
  };

  return (
    <div className="p-4 space-y-2">
      <h2 className="text-xl font-bold">📥 X3DH Receiver (React)</h2>
      <input
        type="text"
        placeholder="EK_A Public (base64)"
        className="border p-2 w-full"
        value={ekAPub}
        onChange={(e) => setEkAPub(e.target.value)}
      />
      <input
        type="text"
        placeholder="IV (base64)"
        className="border p-2 w-full"
        value={iv}
        onChange={(e) => setIv(e.target.value)}
      />
      <textarea
        placeholder="Зашифрованное сообщение (base64)"
        className="border p-2 w-full h-32"
        value={ciphertext}
        onChange={(e) => setCiphertext(e.target.value)}
      />
      <button onClick={decrypt} className="bg-green-600 text-white px-4 py-2 rounded">
        Расшифровать
      </button>
      <div className="border p-2 bg-gray-100 whitespace-pre-wrap">{decrypted}</div>
    </div>
  );
}