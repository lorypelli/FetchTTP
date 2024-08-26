/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    theme: {
        extend: {
            fontFamily: {
                custom: "'Martian Mono', monospace",
            },
            backgroundColor: {
                primary: '#171717',
                secondary: '#232323',
            },
            width: {
                0.25: '1px',
            },
            maxWidth: {
                container: '1000px',
            },
            height: {
                0.25: '1px',
                '1/4': '25vh',
                '1/2': '50vh',
                max: '90vh',
                full: '80vh',
                partial: '70vh',
                fit: '20vh',
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
            },
        },
    },
    plugins: [],
};
