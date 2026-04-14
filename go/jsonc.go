/* Copyright (c) 2021-2025 Richard Rodger, MIT License */

package jsonc

import (
	"regexp"

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

	// Parse grammar text and apply options + val ZZ rule.
	parsed, err := jsonic.Parse(grammarText)
	if err != nil {
		panic("failed to parse jsonc grammar: " + err.Error())
	}
	gm := parsed.(map[string]any)
	optsMap := gm["options"].(map[string]any)

	// @/^\./ resolves to *regexp.Regexp but MapToOptions needs func(string) bool.
	if numMap, ok := optsMap["number"].(map[string]any); ok {
		numMap["exclude"] = regexp.MustCompile(`^\.`).MatchString
	}

	j.Grammar(&jsonic.GrammarSpec{
		OptionsMap: optsMap,
		Rule: map[string]*jsonic.GrammarRuleSpec{
			"val": {
				Open: &jsonic.GrammarAltListSpec{
					Alts:   []*jsonic.GrammarAltSpec{{S: "#ZZ", G: "jsonc"}},
					Inject: &jsonic.GrammarInjectSpec{Append: true},
				},
			},
		},
	})

	// Runtime options not expressible in static grammar.
	j.SetOptions(jsonic.Options{
		Text:    &jsonic.TextOptions{Lex: boolPtr(false)},
		Comment: &jsonic.CommentOptions{Lex: boolPtr(!disallowComments)},
	})
	j.Exclude("jsonic", "imp")

	// Custom value keyword matcher: handles true, false, null.
	// Needed because text lexing is disabled for JSONC compliance
	// (no bare text values), but value keywords must still work.
	VL := j.Token("#VL")
	j.AddMatcher("jsonc-value", 100000, func(lex *jsonic.Lex, rule *jsonic.Rule) *jsonic.Token {
		pnt := lex.Cursor()
		src := lex.Src
		sI := pnt.SI
		if sI >= pnt.Len {
			return nil
		}
		for _, k := range []struct {
			text string
			val  any
		}{{"false", false}, {"true", true}, {"null", nil}} {
			end := sI + len(k.text)
			if end > pnt.Len || src[sI:end] != k.text {
				continue
			}
			if end < pnt.Len {
				ch := src[end]
				if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
					(ch >= '0' && ch <= '9') || ch == '_' || ch == '$' {
					continue
				}
			}
			tkn := lex.Token("#VL", VL, k.val, k.text)
			pnt.SI = end
			pnt.CI += len(k.text)
			return tkn
		}
		return nil
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
