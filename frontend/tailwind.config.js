/** @type {import('tailwindcss').Config} */
export default {
    content: [
        './index.html',
        './src/**/*.{vue,js,ts,jsx,tsx}',
    ],
    theme: {
        extend: {
            height: {
                'scrollbar-full': '73.5vh',
                'scrollbar-partial': '65vh'
            }
        },
    },
    plugins: [],
};