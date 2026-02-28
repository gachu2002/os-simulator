import { forwardRef } from "react";
import type { HTMLAttributes } from "react";

import { cn } from "../../shared/lib/cn";

export const Card = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => {
    return <div ref={ref} className={cn("rounded-xl border border-slate-200 bg-white", className)} {...props} />;
  },
);
Card.displayName = "Card";

export const CardHeader = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => {
    return <div ref={ref} className={cn("p-4 pb-2", className)} {...props} />;
  },
);
CardHeader.displayName = "CardHeader";

export const CardContent = forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => {
    return <div ref={ref} className={cn("p-4 pt-2", className)} {...props} />;
  },
);
CardContent.displayName = "CardContent";

export const CardTitle = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => {
    return <p ref={ref} className={cn("text-base font-semibold text-slate-900", className)} {...props} />;
  },
);
CardTitle.displayName = "CardTitle";

export const CardDescription = forwardRef<HTMLParagraphElement, HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => {
    return <p ref={ref} className={cn("text-sm text-slate-600", className)} {...props} />;
  },
);
CardDescription.displayName = "CardDescription";
