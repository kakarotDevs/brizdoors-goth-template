/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.templ",
    "./layouts/**/*.templ",
    "./components/**/*.templ",
    "./handlers/**/*.go",
    "./*.go",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["'Nunito Sans'", "sans-serif"],
        body: ["'Nunito Sans'", "sans-serif"],
      },
    },
  },
  plugins: [],
};
