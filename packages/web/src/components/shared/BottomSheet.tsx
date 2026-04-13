import type { ReactNode } from "react";
import { Drawer } from "vaul";

interface BottomSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: ReactNode;
}

export function BottomSheet({ open, onOpenChange, children }: BottomSheetProps) {
  return (
    <Drawer.Root open={open} onOpenChange={onOpenChange}>
      <Drawer.Portal>
        <Drawer.Overlay className="fixed inset-0 bg-black/60 z-50" />
        <Drawer.Content className="fixed bottom-0 left-0 right-0 max-w-mobile mx-auto bg-surface-card rounded-t-2xl z-50 p-6">
          <div className="w-12 h-1 bg-surface-secondary rounded-full mx-auto mb-6" />
          {children}
        </Drawer.Content>
      </Drawer.Portal>
    </Drawer.Root>
  );
}
