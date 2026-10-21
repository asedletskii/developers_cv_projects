try {
  var mode = localStorage.getItem('colorMode')
  if (mode !== 'light' && mode !== 'dark') {
    mode = matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
  document.documentElement.dataset.theme = mode
} catch (e) {}
