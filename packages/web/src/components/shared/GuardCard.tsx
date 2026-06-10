import { useNavigate } from "react-router-dom";
import type { LucideIcon } from "lucide-react";
import { Button } from "@/components/ui/button";

interface GuardCardProps {
  icon: LucideIcon;
  title: string;
  message: string;
  actionLabel: string;
  actionPath: string;
}

export function GuardCard({ icon: Icon, title, message, actionLabel, actionPath }: GuardCardProps) {
  const navigate = useNavigate();
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] px-6 text-center">
      <Icon className="mb-4 size-12 text-brand-gold" strokeWidth={1.5} />
      <h2 className="text-lg font-bold text-text-primary mb-2">{title}</h2>
      <p className="text-sm text-text-secondary mb-6">{message}</p>
      <Button onClick={() => navigate(actionPath)} className="w-full max-w-[240px]">
        {actionLabel}
      </Button>
    </div>
  );
}
