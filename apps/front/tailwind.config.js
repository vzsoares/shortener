/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./src/templates/**/*.html", "./src/pages/**/*.html"],
    theme: {
        extend: {
            colors: {
                primary: "var(--primary)",
                background: "var(--background)",
                border: "var(--border)",
                "text-widget": "var(--text-widget)",
                "text-back": "var(--text-back)",
            },
        },
    },
    plugins: [],
};
