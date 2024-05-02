/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    theme: {
        extend: {
            fontFamily: {
                custom: '\'Martian Mono\', monospace',
            },
            backgroundColor: {
                primary: '#171717',
                secondary: '#232323',
            },
            maxWidth: {
                container: '1000px',
            },
            height: {
                0.25: '1px',
                max: '90vh',
                full: '80vh',
                partial: '70vh',
            },
            zIndex: {
                1: 1,
            },
            borderColor: {
                primary: '#171717',
            },
            borderRadius: {
                custom: '10px',
            },
            borderWidth: {
                custom: '5px',
            }
        },
    },
    plugins: [],
};
