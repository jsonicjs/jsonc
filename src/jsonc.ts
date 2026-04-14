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
#
# Function references (@ prefixed) are resolved against the refs map:
#   @exclude-leading-dot  - rejects numbers starting with '.'

{
  options: text: { lex: false }
  options: number: { hex: false oct: false bin: false sep: null exclude: '@exclude-leading-dot' }
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

  // Apply grammar: static options and the val ZZ rule alt.
  const grammar = Jsonic.make()(grammarText)
  jsonic.grammar(grammar)

  // Runtime options that depend on plugin arguments, and
  // number.exclude which requires JS funcref resolution.
  jsonic.options({
    comment: {
      lex: true !== options.disallowComments,
    },
    number: {
      exclude: /^\./,
    },
    rule: {
      include: 'jsonc,json' + (options.allowTrailingComma ? ',comma' : ''),
    },
  })
}

export { Jsonc }

export type { JsoncOptions }
