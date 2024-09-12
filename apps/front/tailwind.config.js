/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./src/pages/**/*.html"],
    theme: {
        extend: {
            colors: {
                primary: "var(--primary)",
                secondary: "var(--secondary)",
                background: "var(--background)",
                border: "var(--border)",
            },
        },
    },
    plugins: [],
};
