import { forwardRef } from "react";
import type { LabelHTMLAttributes } from "react";

import { cn } from "../../shared/lib/cn";

export const Label = forwardRef<HTMLLabelElement, LabelHTMLAttributes<HTMLLabelElement>>(
  ({ className, ...props }, ref) => {
    return <label ref={ref} className={cn("grid gap-1 text-sm text-slate-600", className)} {...props} />;
  },
);

Label.displayName = "Label";
