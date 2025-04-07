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

export default function X3DHSender() {
  const [recipient, setRecipient] = useState<RecipientKeys | null>(null);
  const [message, setMessage] = useState<string>("");
  const [log, setLog] = useState<string[]>([]);

  useEffect(() => {
    fetch("http://localhost:8002/keys")
      .then((res) => res.json())
      .then(setRecipient);
  }, []);

  const encrypt = async () => {
    if (!recipient) return;

    const IK_B = base64ToBytes(recipient.identity_key.public);
    const SPK_B = base64ToBytes(recipient.signed_pre_key.public);
    const OTPK_B = base64ToBytes(recipient.one_time_keys[0].public);
    const IK_SIGN_B = base64ToBytes(recipient.identity_key_signer.public);
    const SPK_SIG = base64ToBytes(recipient.spk_signature);

    const valid = nacl.sign.detached.verify(SPK_B, SPK_SIG, IK_SIGN_B);
    if (!valid) {
      alert("❌ Подпись SPK невалидна!");
      return;
    }

    const EK_A = nacl.box.keyPair();

    const dh = (priv: Uint8Array, pub: Uint8Array) => nacl.scalarMult(priv, pub);
    const dh1 = dh(EK_A.secretKey, IK_B);
    const dh2 = dh(EK_A.secretKey, SPK_B);
    const dh3 = dh(EK_A.secretKey, OTPK_B);

    const combined = new Uint8Array([...dh1, ...dh2, ...dh3]);
    const hashBuffer = await crypto.subtle.digest("SHA-256", combined);
    const sharedSecret = new Uint8Array(hashBuffer);

    const iv = crypto.getRandomValues(new Uint8Array(12));
    const key = await crypto.subtle.importKey("raw", sharedSecret, { name: "AES-GCM" }, false, ["encrypt"]);
    const encoded = new TextEncoder().encode(message);
    const ciphertext = await crypto.subtle.encrypt({ name: "AES-GCM", iv }, key, encoded);

    setLog((log) => [
      `🔒 Encrypted: ${bytesToBase64(new Uint8Array(ciphertext))}`,
      `📎 IV: ${bytesToBase64(iv)}`,
      `📤 EK_A Public: ${bytesToBase64(EK_A.publicKey)}`,
      ...log,
    ]);
  };

  return (
    <div className="p-4 space-y-2">
      <h2 className="text-xl font-bold">🔐 X3DH Sender (React)</h2>
      <input
        type="text"
        className="border p-2 w-full"
        placeholder="Введите сообщение"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button onClick={encrypt} className="bg-blue-500 text-white px-4 py-2 rounded">
        Зашифровать и отправить
      </button>
      <textarea className="w-full h-60 p-2 border" readOnly value={log.join("\n")} />
    </div>
  );
}