import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { ArrowLeft, CheckCircle2, ExternalLink, Loader2, Wallet } from "lucide-react";
import { QRCodeSVG } from "qrcode.react";

const RPC = "https://api.devnet.solana.com";
const USDC_MINT = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU";
const RECEIVER = "BEuTBCN3EygdmxRQXcWuYYXAUcwXtVnPbPSTvsgDz5Lc";
const AMOUNT_USDC = 10;
const TOKEN_UNITS = 10_000_000n;

type RpcResponse<T> = { result?: T; error?: { message?: string } };
type SignatureInfo = { signature: string; err: unknown | null };
type ParsedInstruction = { program: string; parsed?: { type?: string; info?: Record<string, unknown> } };
type ParsedTransaction = {
  meta?: { err: unknown | null };
  transaction?: { message?: { instructions?: ParsedInstruction[] } };
};
type ParsedTokenAccount = { value: null | { data?: { parsed?: { info?: { mint?: string; owner?: string } } } } };

function base58(bytes: Uint8Array): string {
  const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";
  let value = 0n;
  for (const byte of bytes) value = (value << 8n) | BigInt(byte);
  let encoded = "";
  while (value > 0n) {
    const remainder = Number(value % 58n);
    value /= 58n;
    encoded = alphabet[remainder] + encoded;
  }
  for (const byte of bytes) {
    if (byte !== 0) break;
    encoded = "1" + encoded;
  }
  return encoded || "1";
}

function makeReference(): string {
  const bytes = new Uint8Array(32);
  crypto.getRandomValues(bytes);
  return base58(bytes);
}

async function rpc<T>(method: string, params: unknown[]): Promise<T> {
  const response = await fetch(RPC, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ jsonrpc: "2.0", id: 1, method, params }),
  });
  if (!response.ok) throw new Error("Solana Devnet RPC is temporarily unavailable.");
  const body = (await response.json()) as RpcResponse<T>;
  if (body.error) throw new Error(body.error.message || "Devnet RPC request failed.");
  if (body.result === undefined) throw new Error("Devnet RPC returned no result.");
  return body.result;
}

export default function SolanaCheckout() {
  const [reference] = useState(makeReference);
  const [confirmedSignature, setConfirmedSignature] = useState("");
  const [error, setError] = useState("");
  const [checking, setChecking] = useState(true);

  const paymentUrl = useMemo(() => {
    const query = new URLSearchParams({
      amount: String(AMOUNT_USDC),
      "spl-token": USDC_MINT,
      reference,
      label: "LinkWorld",
      message: "10 USDC Devnet eSIM plan · demo only",
    });
    return `solana:${RECEIVER}?${query.toString()}`;
  }, [reference]);

  useEffect(() => {
    let active = true;
    let timer: ReturnType<typeof setTimeout> | undefined;
    const seen = new Set<string>();

    const check = async () => {
      try {
        const signatures = await rpc<SignatureInfo[]>("getSignaturesForAddress", [reference, { limit: 10 }]);
        for (const item of signatures) {
          if (seen.has(item.signature)) continue;
          seen.add(item.signature);
          if (item.err) continue;

          const tx = await rpc<ParsedTransaction | null>("getTransaction", [
            item.signature,
            { encoding: "jsonParsed", commitment: "confirmed", maxSupportedTransactionVersion: 0 },
          ]);
          if (!tx || tx.meta?.err) continue;
          const instructions = tx.transaction?.message?.instructions ?? [];

          for (const ix of instructions) {
            if (ix.program !== "spl-token" || ix.parsed?.type !== "transferChecked") continue;
            const info = ix.parsed.info ?? {};
            const amount = info.tokenAmount as { amount?: string; decimals?: number } | undefined;
            const destination = info.destination;
            if (
              info.mint !== USDC_MINT ||
              amount?.amount !== String(TOKEN_UNITS) ||
              amount?.decimals !== 6 ||
              typeof destination !== "string"
            ) continue;

            const account = await rpc<ParsedTokenAccount>("getAccountInfo", [
              destination,
              { encoding: "jsonParsed", commitment: "confirmed" },
            ]);
            const accountInfo = account.value?.data?.parsed?.info;
            if (accountInfo?.mint === USDC_MINT && accountInfo.owner === RECEIVER) {
              if (active) {
                setError("");
                setConfirmedSignature(item.signature);
              }
              return;
            }
          }
        }
        if (active) timer = setTimeout(check, 3000);
      } catch (cause) {
        if (active) {
          setError(cause instanceof Error ? cause.message : "Could not check the payment.");
          timer = setTimeout(check, 8000);
        }
      }
    };

    void check();
    return () => {
      active = false;
      if (timer) clearTimeout(timer);
    };
  }, [reference]);

  useEffect(() => {
    if (confirmedSignature) setChecking(false);
  }, [confirmedSignature]);

  return (
    <main className="min-h-screen bg-surface text-text-primary px-5 py-6 max-w-mobile mx-auto">
      <Link to="/" className="inline-flex items-center gap-2 text-sm text-text-secondary hover:text-brand-gold">
        <ArrowLeft className="size-4" /> Back
      </Link>
      <section className="mt-10 rounded-2xl bg-card p-6 text-center shadow-lg">
        <div className="mx-auto mb-5 grid size-14 place-items-center rounded-full bg-brand-royal/20 text-brand-gold">
          {confirmedSignature ? <CheckCircle2 className="size-7" /> : <Wallet className="size-7" />}
        </div>
        <h1 className="text-2xl font-bold">Solana Devnet eSIM Checkout</h1>
        <p className="mt-3 text-sm text-text-secondary">
          Pay 10 test USDC on Solana Devnet. This is a demo checkout; no real eSIM service is activated.
        </p>
        <div className="mt-6 flex justify-center rounded-xl bg-white p-4">
          <QRCodeSVG value={paymentUrl} size={220} level="M" />
        </div>
        <a href={paymentUrl} className="mt-5 inline-flex items-center justify-center gap-2 rounded-xl bg-brand-gold px-5 py-3 font-semibold text-surface">
          Open in Solana wallet <ExternalLink className="size-4" />
        </a>
        <p className="mt-4 break-all text-xs text-text-secondary">Recipient: {RECEIVER}</p>
        <p className="mt-1 text-xs text-text-muted">Token: Devnet USDC · Mint {USDC_MINT}</p>
        <div className="mt-6 min-h-6 text-sm" aria-live="polite">
          {confirmedSignature ? (
            <div className="text-status-success">
              <p>Payment confirmed. Demo order is marked paid.</p>
              <a className="underline" target="_blank" rel="noreferrer" href={`https://explorer.solana.com/tx/${confirmedSignature}?cluster=devnet`}>View Devnet transaction</a>
            </div>
          ) : checking ? (
            <span className="inline-flex items-center gap-2 text-text-secondary"><Loader2 className="size-4 animate-spin" /> Waiting for a confirmed USDC payment…</span>
          ) : null}
          {error && <p className="mt-2 text-status-danger">{error}</p>}
        </div>
        <p className="mt-5 text-xs text-text-muted">
          Use a Solana wallet with Devnet enabled. You need Devnet SOL for fees and Devnet USDC for payment.
        </p>
        <div className="mt-2 flex justify-center gap-4 text-xs">
          <a className="underline" href="https://faucet.solana.com/" target="_blank" rel="noreferrer">SOL faucet</a>
          <a className="underline" href="https://faucet.circle.com/" target="_blank" rel="noreferrer">USDC faucet</a>
        </div>
      </section>
    </main>
  );
}
