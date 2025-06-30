/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.templ",
    "./layouts/**/*.templ",
    "./components/**/*.templ",
    "./handlers/**/*.go",
    "./*.go",
  ],
  safelist: [
    'bg-bg-dark',
    'bg-bg-light', 
    'border-bg-light'
  ],
  darkMode: 'class', // Enable class-based dark mode
  theme: {
    extend: {
      fontFamily: {
        sans: ["'Nunito Sans'", "sans-serif"],
        body: ["'Nunito Sans'", "sans-serif"],
      },
      colors: {
        // Custom HSL-based color system
        // Dark mode colors (default)
        'bg-dark': 'hsl(var(--bg-dark))',
        'bg': 'hsl(var(--bg))',
        'bg-light': 'hsl(var(--bg-light))',
        'text': 'hsl(var(--text))',
        'text-muted': 'hsl(var(--text-muted))',
      },
    },
  },
  plugins: [],
};
