/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.templ",
    "./layouts/**/*.templ",
    "./components/**/*.templ",
    "./handlers/**/*.go",
    "./*.go",
  ],
  safelist: ["bg-bg-dark", "bg-bg-light", "border-bg-light"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["'Nunito Sans'", "system-ui", "-apple-system", "sans-serif"],
      },
      colors: {
        // Custom HSL-based color system for light theme
        // bg-bg-dark is a slightly darker light shade for hover/active states
        "bg-dark": "hsl(var(--bg-dark))",
        bg: "hsl(var(--bg))",
        "bg-light": "hsl(var(--bg-light))",
        text: "hsl(var(--text))",
        "text-muted": "hsl(var(--text-muted))",
        caret: "hsl(var(--text))", // Caret color matching text color
      },
      letterSpacing: {
        apple: "-0.01em",
        "apple-tight": "-0.02em",
      },
    },
  },
  plugins: [],
};
