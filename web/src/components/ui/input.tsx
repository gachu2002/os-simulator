import { forwardRef } from "react";
import type { InputHTMLAttributes } from "react";

import { cn } from "../../shared/lib/cn";

export const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => {
    return (
      <input
        className={cn(
          "h-9 rounded-md border border-slate-300 bg-white px-2.5 py-1.5 text-sm text-slate-900",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400",
          className,
        )}
        ref={ref}
        {...props}
      />
    );
  },
);

Input.displayName = "Input";
