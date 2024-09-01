/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{svelte,js,ts,jsx,tsx}'],
    theme: {
        extend: {
            backgroundColor: {
                primary: '#171717',
                secondary: '#232323',
            },
            fontFamily: {
                custom: "'Martian Mono', monospace",
            },
            height: {
                max: '90vh',
                full: '85vh',
            },
        },
    },
    plugins: [],
};
