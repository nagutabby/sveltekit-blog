import daisyui from 'daisyui'
import typography from '@tailwindcss/typography'

export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {},
  },
  plugins: [
    daisyui, typography
  ],
  daisyui: {
    themes: [
      "dracula"
    ],
  },
}

