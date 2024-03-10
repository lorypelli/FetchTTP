/** @type {import('tailwindcss').Config} */
export default {
    content: [
        './index.html',
        './src/**/*.{vue,js,ts,jsx,tsx}',
    ],
    theme: {
        extend: {
            height: {
                'scrollbar-full': '83.5vh',
                'scrollbar-partial': '75vh'
            }
        },
    },
    plugins: [],
};