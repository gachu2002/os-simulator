import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      borderRadius: {
        lg: "0.75rem",
        md: "0.625rem",
        sm: "0.5rem",
      },
      colors: {
        border: "hsl(34 19% 82%)",
        background: "hsl(42 24% 94%)",
        foreground: "hsl(214 12% 14%)",
        card: "hsl(40 33% 98%)",
        "card-foreground": "hsl(214 12% 14%)",
        primary: "hsl(202 60% 30%)",
        "primary-foreground": "hsl(0 0% 100%)",
        muted: "hsl(210 10% 40%)",
      },
    },
  },
  plugins: [],
};

export default config;
