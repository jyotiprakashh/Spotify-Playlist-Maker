/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx,mdx}"],
  theme: {
    extend: {
      screens: {
        sm: "480px",
        md: "768px",
        lg: "976px",
        xl: "1440px",
      },

      colors: {
        spotifygreen: "#1db954",
        spotifyblack: "#191414",
        spotifywhite: "#ffffff",
        spotifygray: "#a0a0a0",
        spotifydarkgray: "#282828",
        spotifylightgray: "#b3b3b3",
      },
      animation: {
        'fade-in': 'fadeIn 1s ease-in forwards',
        'bounce': 'bounce 1s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: 0 },
          '100%': { opacity: 1 },
        },
      },
    },
  },
  plugins: [],
};
