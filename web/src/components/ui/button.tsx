import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import type { ButtonHTMLAttributes } from "react";

import { cn } from "../../shared/lib/cn";

const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-md border font-semibold transition disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "border-blue-800 bg-blue-700 text-white hover:bg-blue-800",
        secondary: "border-slate-300 bg-slate-200 text-slate-900 hover:bg-slate-300",
        outline: "border-slate-400 bg-white text-slate-900 hover:bg-slate-50",
        destructive: "border-red-800 bg-red-700 text-white hover:bg-red-800",
        success: "border-emerald-800 bg-emerald-700 text-white hover:bg-emerald-800",
      },
      size: {
        default: "h-9 px-3 py-1.5 text-sm",
        sm: "h-8 px-2.5 text-xs",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

export function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: ButtonProps) {
  const Comp = asChild ? Slot : "button";
  return <Comp className={cn(buttonVariants({ variant, size }), className)} {...props} />;
}
