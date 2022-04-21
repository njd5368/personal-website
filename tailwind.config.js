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
    },
  },
  plugins: [],
}
