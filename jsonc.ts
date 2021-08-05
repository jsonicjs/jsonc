
import { Jsonic } from 'jsonic'


type JsoncOptions = {
  allowTrailingComma?: boolean
  disallowComments?: boolean
}

function Jsonc(jsonic: Jsonic, options: JsoncOptions) {
  jsonic.options({
    // TODO: replace with a Jsonic.options('json') preset
    // TODO: need to accept params for rule include ... hmmm
    text: { lex: false },
    number: { hex: false, oct: false, bin: false, sep: null },
    string: { chars: '"', multiChars: '', allowUnknown: false, escape: { v: null } },
    comment: { lex: false },
    map: { extend: false },
    rule: {
      finish: false,
      include: 'json' + (options.allowTrailingComma ? ',comma' : '')
    },
  })

  jsonic.options({
    comment: {
      lex: true && !options.disallowComments,
      marker: [
        { line: true, start: '//', lex: true },
        { line: false, start: '/' + '*', end: '*' + '/', lex: true },
      ],
    }
  })
}

export {
  Jsonc,
}

export type {
  JsoncOptions,
}
