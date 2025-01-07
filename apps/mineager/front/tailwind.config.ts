import { Config } from 'tailwindcss';

export default {
  content: ['./front/public/index.html', './front/src/**/*.{js,ts,html}'],
  theme: {
    extend: {
      colors: {
        primary: '#00B1CC',
        background: '#141b33',
        'background-darker': '#0b1021',
        'background-hover': '#060913',
        info: '#04e4f4',
        success: '#21FA00',
        warning: '#f6f000',
        error: '#fca016',
        danger: '#f80363',
        neutral: '#272f4d',
      },
      screens: {
        phone: { max: '599px' },
        tablet: { min: '600px', max: '1149px' },
        laptop: { min: '1150px', max: '1799px' },
        desktop: { min: '1800px' },
        tv: { min: '2800px' },
      },
      minWidth: {
        '1rem': '1rem',
        '1.5rem': '1.5rem',
        '2rem': '2rem',
        '3rem': '3rem',
      },
      maxWidth: {
        '1rem': '1rem',
        '1.5rem': '1.5rem',
        '2rem': '2rem',
        '3rem': '3rem',
      },
      minHeight: {
        '1rem': '1rem',
        '1.5rem': '1.5rem',
        '2rem': '2rem',
        '3rem': '3rem',
      },
      maxHeight: {
        '1rem': '1rem',
        '1.5rem': '1.5rem',
        '2rem': '2rem',
        '3rem': '3rem',
      },
      fontSize: {
        '1rem': '1rem',
        '1.5rem': '1.5rem',
        '2rem': '2rem',
        '2.5rem': '2.5rem',
      },
      height: {
        '15': '3.75rem',
      },
    },
  },
  plugins: [],
} satisfies Config;
