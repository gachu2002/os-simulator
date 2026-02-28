import { cva, type VariantProps } from "class-variance-authority";
import type { HTMLAttributes } from "react";

import { cn } from "../../shared/lib/cn";

const badgeVariants = cva(
  "inline-flex rounded-full px-2 py-0.5 text-[11px] font-semibold uppercase tracking-wide",
  {
    variants: {
      variant: {
        default: "bg-slate-200 text-slate-700",
        secondary: "bg-blue-100 text-blue-700",
        destructive: "bg-red-100 text-red-700",
        success: "bg-emerald-100 text-emerald-700",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);

interface BadgeProps extends HTMLAttributes<HTMLSpanElement>, VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}
