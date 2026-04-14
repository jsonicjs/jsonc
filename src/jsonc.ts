/* Copyright (c) 2021-2025 Richard Rodger, MIT License */

// Import Jsonic types used by plugin.
import { Jsonic } from 'jsonic'

type JsoncOptions = {
  allowTrailingComma?: boolean
  disallowComments?: boolean
}

// --- BEGIN EMBEDDED jsonc-grammar.jsonic ---
const grammarText = `
# JSONC Grammar Definition
# Parsed by a standard Jsonic instance and passed to jsonic.grammar()
# Extends standard JSON grammar with end-of-input value handling.
# Trailing commas are added programmatically via rule modification.

{
  options: text: { lex: false }
  options: number: { hex: false oct: false bin: false sep: null exclude: "@/^\\\\./" }
  options: string: { chars: '"' multiChars: '' allowUnknown: false }
  options: string: escape: { v: null }
  options: map: { extend: false }
  options: lex: { empty: false }
  options: rule: { finish: false }

  rule: val: open: {
    alts: [
      { s: '#ZZ' g: jsonc }
    ]
    inject: { append: true }
  }
}
`
// --- END EMBEDDED jsonc-grammar.jsonic ---

function Jsonc(jsonic: Jsonic, options: JsoncOptions) {

  // Apply grammar: static options and val ZZ rule alt.
  jsonic.grammar(Jsonic.make()(grammarText))

  // Runtime options that depend on plugin arguments.
  jsonic.options({
    comment: {
      lex: true !== options.disallowComments,
    },
    rule: {
      include: 'jsonc,json' + (options.allowTrailingComma ? ',comma' : ''),
    },
  })
}

export { Jsonc }

export type { JsoncOptions }
