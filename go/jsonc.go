/* Copyright (c) 2021-2025 Richard Rodger, MIT License */

package jsonc

import (
	jsonic "github.com/jsonicjs/jsonic/go"
)

// JsoncOptions configures the JSONC parser.
type JsoncOptions struct {
	// AllowTrailingComma enables trailing commas in objects and arrays.
	// Default: false (standard JSONC behavior; set true for VS Code compatibility).
	AllowTrailingComma *bool
	// DisallowComments disables comment parsing.
	// Default: false (comments are enabled by default in JSONC).
	DisallowComments *bool
}

// Parse parses a JSONC string and returns the result.
// Returns the parsed value (map, slice, string, float64, bool, or nil) and any error.
func Parse(src string, opts ...JsoncOptions) (any, error) {
	var o JsoncOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	j := MakeJsonic(o)
	return j.Parse(src)
}

// MakeJsonic creates a jsonic instance configured for JSONC parsing.
func MakeJsonic(opts ...JsoncOptions) *jsonic.Jsonic {
	var o JsoncOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	// lex.empty is stored on the Jsonic struct in Make(), not in the config.
	j := jsonic.Make(jsonic.Options{
		Lex: &jsonic.LexOptions{Empty: boolPtr(false)},
	})

	j.Use(jsoncPlugin, optionsToMap(&o))

	return j
}

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

// jsoncPlugin is the jsonic plugin that configures JSONC parsing.
func jsoncPlugin(j *jsonic.Jsonic, pluginOpts map[string]any) {
	allowTrailingComma, _ := pluginOpts["allowTrailingComma"].(bool)
	disallowComments, _ := pluginOpts["disallowComments"].(bool)

	// Apply grammar options from text.
	if err := j.GrammarText(grammarText); err != nil {
		panic("failed to apply jsonc grammar: " + err.Error())
	}

	// Apply val ZZ rule (GrammarText handles options only, not rules).
	j.Grammar(&jsonic.GrammarSpec{
		Rule: map[string]*jsonic.GrammarRuleSpec{
			"val": {
				Open: &jsonic.GrammarAltListSpec{
					Alts:   []*jsonic.GrammarAltSpec{{S: "#ZZ", G: "jsonc"}},
					Inject: &jsonic.GrammarInjectSpec{Append: true},
				},
			},
		},
	})

	// Runtime options and options not handled by MapToOptions (text).
	j.SetOptions(jsonic.Options{
		Text:    &jsonic.TextOptions{Lex: boolPtr(false)},
		Comment: &jsonic.CommentOptions{Lex: boolPtr(!disallowComments)},
		Rule:    &jsonic.RuleOptions{Exclude: "jsonic,imp"},
	})

	// Trailing comma support.
	if allowTrailingComma {
		CA, CB, CS := j.Token("#CA"), j.Token("#CB"), j.Token("#CS")
		j.Rule("pair", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{S: [][]jsonic.Tin{{CA}, {CB}}, B: 1})
		})
		j.Rule("elem", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{S: [][]jsonic.Tin{{CA}, {CS}}, B: 1})
		})
	}
}

func optionsToMap(o *JsoncOptions) map[string]any {
	return map[string]any{
		"allowTrailingComma": boolOpt(o.AllowTrailingComma, false),
		"disallowComments":   boolOpt(o.DisallowComments, false),
	}
}

func boolOpt(p *bool, def bool) bool {
	if p != nil {
		return *p
	}
	return def
}

func boolPtr(b bool) *bool {
	return &b
}
