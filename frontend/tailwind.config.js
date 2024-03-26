/** @type {import('tailwindcss').Config} */
export default {
    content: [
        './index.html',
        './src/**/*.{vue,js,ts,jsx,tsx}',
    ],
    theme: {
        extend: {
            height: {
                'max': '100vh',
                'full': '80vh'
            }
        },
    },
    plugins: [],
};