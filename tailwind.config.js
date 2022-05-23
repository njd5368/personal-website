module.exports = {
  content: [
      './site/templates/**/*.gohtml'
  ],
  theme: {
    fontFamily: {
      'sans': ['Montserrat']
    },
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      'black': '#000',
      'white': '#fff',
      'gray': {
        100: '#D9E0EE',
        200: '#C3BAC6',
        300: '#988BA2',
        400: '#6E6C7E',
        500: '#575268',
        600: '#302D41',
        700: '#1E1E2E',
        800: '#1A1826',
        900: '#161320',
      },
      'flamingo': '#F2CDCD',
      'mauve':    '#DDB6F2',
      'pink':     '#F5C2E7',
      'maroon':   '#E8A2AF',
      'red':      '#F28FAD',
      'peach':    '#F8BD96',
      'yellow':   '#FAE3B0',
      'green':    '#ABE9B3',
      'teal':     '#B5E8E0',
      'blue':     '#96CDFB',
      'sky':      '#89DCEB',
    },
    extend: {
      spacing: {
        '26': '6.5rem',
      },
      aspectRatio: {
        'golden-y': '1 / 1.618',
        'golden-x': '1.618 / 1',
      },
      typography: (theme) => ({
        catppuccin: {
          css: {
            '--tw-prose-body': theme('colors.gray[100]'),
            '--tw-prose-headings': theme('colors.gray[100]'),
            '--tw-prose-lead': theme('colors.gray[100]'),
            '--tw-prose-links': theme('colors.blue'),
            '--tw-prose-bold': theme('colors.white'),
            '--tw-prose-counters': theme('colors.gray[100]'),
            '--tw-prose-bullets': theme('colors.gray[100]'),
            '--tw-prose-hr': theme('colors.gray[100]'),
            '--tw-prose-quotes': theme('colors.gray[100]'),
            '--tw-prose-quote-borders': theme('colors.gray[100]'),
            '--tw-prose-captions': theme('colors.gray[100]'),
            '--tw-prose-code': theme('colors.gray[100]'),
            '--tw-prose-pre-code': theme('colors.gray[100]'),
            '--tw-prose-pre-bg': theme('colors.gray[100]'),
            '--tw-prose-th-borders': theme('colors.gray[100]'),
            '--tw-prose-td-borders': theme('colors.gray[100]'),
            '--tw-prose-invert-body': theme('colors.gray[100]'),
            '--tw-prose-invert-headings': theme('colors.white'),
            '--tw-prose-invert-lead': theme('colors.gray[100]'),
            '--tw-prose-invert-links': theme('colors.white'),
            '--tw-prose-invert-bold': theme('colors.white'),
            '--tw-prose-invert-counters': theme('colors.gray[100]'),
            '--tw-prose-invert-bullets': theme('colors.gray[100]'),
            '--tw-prose-invert-hr': theme('colors.gray[100]'),
            '--tw-prose-invert-quotes': theme('colors.gray[100]'),
            '--tw-prose-invert-quote-borders': theme('colors.gray[100]'),
            '--tw-prose-invert-captions': theme('colors.gray[100]'),
            '--tw-prose-invert-code': theme('colors.gray[100]'),
            '--tw-prose-invert-pre-code': theme('colors.gray[100]'),
            '--tw-prose-invert-pre-bg': 'rgb(0 0 0 / 50%)',
            '--tw-prose-invert-th-borders': theme('colors.gray[100]'),
            '--tw-prose-invert-td-borders': theme('colors.gray[100]'),
          },
        },
      })

    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
