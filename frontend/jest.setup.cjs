// jsdom не отдаёт TextEncoder/TextDecoder, а react-router обращается к ним при импорте
const { TextDecoder, TextEncoder } = require('util')

Object.assign(globalThis, { TextDecoder, TextEncoder })

// jsdom не реализует matchMedia, а gsap/ScrollTrigger дёргает его при регистрации плагина
// (что происходит уже при импорте `@shared/animations`) — без полифилла падает любой тест,
// который транзитивно доходит до компонентов страницы резюме.
window.matchMedia =
  window.matchMedia ||
  function matchMedia() {
    return {
      matches: false,
      media: '',
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    }
  }
