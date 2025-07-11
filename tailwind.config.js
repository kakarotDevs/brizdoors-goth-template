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
      animation: {
        "spin-slow": "spin 2s linear infinite",
        "door-swing": "doorSwing 1.5s ease-in-out infinite",
        "door-swing-slow": "doorSwing 2.5s ease-in-out infinite",
      },
      keyframes: {
        doorSwing: {
          "0%, 100%": {
            transform: "skewY(0deg) scaleX(1)",
            filter: "brightness(1)",
          },
          "50%": {
            transform: "skewY(35deg) scaleX(0.6)",
            filter: "brightness(0.8)",
          },
        },
      },
    },
  },
  plugins: [
    function ({ addUtilities }) {
      const newUtilities = {
        ".perspective-1000": {
          perspective: "1000px",
        },
        ".preserve-3d": {
          "transform-style": "preserve-3d",
        },
      };
      addUtilities(newUtilities);
    },
  ],
};
