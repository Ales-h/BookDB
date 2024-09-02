/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
    './public/**/*.html',
    './public/**/*.js',
    ],
  theme: {
        extend: {
            colors: {
                primary: '#3490dc',  // Custom primary color
                secondary: '#ffed4a', // Custom secondary color
                boxcolor: '#1a202c',   // Dark background color
                bgcolor: '#11151f',
                lighttext: '#f7fafc', // Light text color
        },
            fontFamily: {
                sans: ['Arial', 'Helvetica', 'sans-serif'],
                serif: ['Georgia', 'serif'],
         },
        },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    ],
}

