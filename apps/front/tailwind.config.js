/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./src/pages/**/*.html"],
    theme: {
        extend: {
            colors: {
                primary: "#532B88",
                secondary: "#2F184B",
                neutral: {
                    500: "#F4EFFA",
                },
                light: "#9B72CF",
                lighter: "#C8B1E4",
            },
        },
    },
    plugins: [],
};
