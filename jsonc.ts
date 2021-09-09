import { Jsonic } from 'jsonic'

type JsoncOptions = {
  allowTrailingComma?: boolean
  disallowComments?: boolean
}

function Jsonc(jsonic: Jsonic, options: JsoncOptions) {
  jsonic.options({
    text: {
      lex: false,
    },
    number: {
      hex: false,
      oct: false,
      bin: false,
      sep: null,
    },
    string: {
      chars: '"',
      multiChars: '',
      allowUnknown: false,
      escape: {
        v: null,
      },
    },
    comment: {
      lex: true !== options.disallowComments,
    },
    map: {
      extend: false,
    },
    lex: {
      empty: false,
    },
    rule: {
      finish: false,
      include: 'json' + (options.allowTrailingComma ? ',comma' : ''),
    },
  })
}

export { Jsonc }

export type { JsoncOptions }
