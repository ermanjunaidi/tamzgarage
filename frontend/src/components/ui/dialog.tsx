import * as React from "react"
import { cn } from "@/lib/utils"

function Dialog({ open, onOpenChange, children }: { open: boolean; onOpenChange: (open: boolean) => void; children: React.ReactNode }) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/40" onClick={() => onOpenChange(false)} />
      <div className="relative z-50 w-full max-w-lg max-h-[85vh] overflow-auto rounded-xl border bg-background shadow-lg mx-4">
        {children}
      </div>
    </div>
  )
}

function DialogContent({ className, children, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("p-6", className)} {...props}>{children}</div>
}

function DialogHeader({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("mb-4", className)} {...props} />
}

function DialogTitle({ className, ...props }: React.ComponentProps<"h2">) {
  return <h2 className={cn("text-lg font-semibold", className)} {...props} />
}

export { Dialog, DialogContent, DialogHeader, DialogTitle }
