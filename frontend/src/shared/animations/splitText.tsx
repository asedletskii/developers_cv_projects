import { Fragment } from 'react'

export const splitText = (text: string, className: string) =>
  text.split(' ').map((word, wordIndex, words) => (
    <Fragment key={wordIndex}>
      <span style={{ display: 'inline-block', whiteSpace: 'nowrap' }}>
        {word.split('').map((char, index) => (
          <span key={index} style={{ display: 'inline-block' }} className={className}>
            {char}
          </span>
        ))}
      </span>
      {wordIndex < words.length - 1 && ' '}
    </Fragment>
  ))
