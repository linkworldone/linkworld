import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAccount } from "wagmi";
import { ConnectButton } from "@/components/wallet/ConnectButton";
import { RegisterSheet } from "@/components/wallet/RegisterSheet";
import { useUser } from "@/hooks/useUser";

export default function Landing() {
  const navigate = useNavigate();
  const { address, isConnected } = useAccount();
  const { data: user, isLoading } = useUser(address);
  const [showRegister, setShowRegister] = useState(false);

  useEffect(() => {
    if (!isConnected || isLoading) return;
    if (user) {
      navigate("/dashboard", { replace: true });
    } else if (isConnected && !user) {
      setShowRegister(true);
    }
  }, [isConnected, user, isLoading, navigate]);

  return (
    <div className="min-h-screen bg-surface text-text-primary flex flex-col max-w-mobile mx-auto">
      <div className="px-5 py-4 flex justify-between items-center">
        <span className="font-orbitron text-lg font-extrabold bg-gradient-to-r from-brand-blue to-brand-purple bg-clip-text text-transparent">
          LinkWorld
        </span>
        <ConnectButton />
      </div>

      <div className="flex-1 flex flex-col items-center justify-center px-6 text-center">
        <div className="w-[72px] h-[72px] rounded-full bg-gradient-to-br from-brand-blue to-brand-purple flex items-center justify-center text-3xl mb-7">
          🌐
        </div>
        <h1 className="text-[32px] font-extrabold leading-tight mb-3">
          Global Seamless<br />Communication
        </h1>
        <p className="text-sm text-text-secondary leading-relaxed mb-9">
          Connect your wallet. Deposit.<br />
          Get a virtual number.<br />
          No KYC. Pay after use.
        </p>
        <ConnectButton label="Get Started" />
        <button className="mt-3 text-sm text-text-secondary border border-border rounded-xl px-12 py-3 w-full max-w-[280px]">
          Learn More
        </button>
      </div>

      <div className="flex justify-around py-6 border-t border-border">
        <div className="text-center">
          <div className="text-xl font-extrabold text-brand-blue">50+</div>
          <div className="text-[11px] text-text-muted mt-0.5">Countries</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-extrabold text-status-success">2.5%</div>
          <div className="text-[11px] text-text-muted mt-0.5">Platform Fee</div>
        </div>
        <div className="text-center">
          <div className="text-xl font-extrabold text-brand-purple">0</div>
          <div className="text-[11px] text-text-muted mt-0.5">KYC Required</div>
        </div>
      </div>

      <div className="text-center py-3 text-[10px] text-text-muted">
        Powered by 0G Chain & Chainlink
      </div>

      {address && (
        <RegisterSheet
          address={address}
          open={showRegister}
          onClose={() => setShowRegister(false)}
          onSuccess={() => navigate("/dashboard", { replace: true })}
        />
      )}
    </div>
  );
}
