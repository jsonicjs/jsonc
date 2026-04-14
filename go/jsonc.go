/* Copyright (c) 2021-2025 Richard Rodger, MIT License */

package jsonc

import (
	jsonic "github.com/jsonicjs/jsonic/go"
)

// --- BEGIN EMBEDDED jsonc-grammar.jsonic ---
const grammarText = `
# JSONC Grammar Definition
# Parsed by a standard Jsonic instance and passed to jsonic.grammar()
# Extends standard JSON grammar with end-of-input value handling.
# Trailing commas are added programmatically via rule modification.

{
  options: text: { lex: false }
  options: number: { hex: false oct: false bin: false sep: null exclude: "@/^\\./" }
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

// Jsonc is the jsonic plugin that configures JSONC parsing.
func Jsonc(j *jsonic.Jsonic, pluginOpts map[string]any) {
	commentLex := true != toBool(pluginOpts["disallowComments"])
	ruleExclude := "jsonic,imp,comma"
	if toBool(pluginOpts["allowTrailingComma"]) {
		ruleExclude = "jsonic,imp"
	}

	// Apply grammar: static options and val ZZ rule alt.
	if err := j.GrammarText(grammarText); err != nil {
		panic("failed to apply jsonc grammar: " + err.Error())
	}

	// Runtime options that depend on plugin arguments.
	j.SetOptions(jsonic.Options{
		Comment: &jsonic.CommentOptions{Lex: &commentLex},
		Rule:    &jsonic.RuleOptions{Exclude: ruleExclude},
	})

	// Trailing comma support (Go jsonic has no built-in "comma" group alts).
	if toBool(pluginOpts["allowTrailingComma"]) {
		CA, CB, CS := j.Token("#CA"), j.Token("#CB"), j.Token("#CS")
		j.Rule("pair", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{S: [][]jsonic.Tin{{CA}, {CB}}, B: 1})
		})
		j.Rule("elem", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{S: [][]jsonic.Tin{{CA}, {CS}}, B: 1})
		})
	}
}

func toBool(v any) bool {
	b, _ := v.(bool)
	return b
}
