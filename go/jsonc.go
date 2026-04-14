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

	disallowComments := boolOpt(o.DisallowComments, false)
	commentLex := !disallowComments

	jopts := jsonic.Options{
		Text: &jsonic.TextOptions{
			Lex: boolPtr(false),
		},
		Number: &jsonic.NumberOptions{
			Lex: boolPtr(true),
			Hex: boolPtr(false),
			Oct: boolPtr(false),
			Bin: boolPtr(false),
			Exclude: func(s string) bool {
				return len(s) > 0 && s[0] == '.'
			},
		},
		String: &jsonic.StringOptions{
			Chars:        `"`,
			AllowUnknown: boolPtr(false),
		},
		Comment: &jsonic.CommentOptions{
			Lex: &commentLex,
		},
		Map: &jsonic.MapOptions{
			Extend: boolPtr(false),
		},
		Rule: &jsonic.RuleOptions{
			Finish:  boolPtr(false),
			Exclude: "jsonic,imp",
		},
		Lex: &jsonic.LexOptions{
			Empty: boolPtr(false),
		},
	}

	j := jsonic.Make(jopts)

	pluginMap := optionsToMap(&o)
	j.Use(jsoncPlugin, pluginMap)

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
  options: number: { hex: false oct: false bin: false sep: null }
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

	VL := j.Token("#VL")
	ZZ := j.Token("#ZZ")

	// Custom value keyword matcher: handles true, false, null.
	// Needed because text lexing is disabled for JSONC compliance
	// (no bare text values allowed), but value keywords must still work.
	// Priority 100000 runs before all built-in matchers (same pattern as ini plugin).
	j.AddMatcher("jsonc-value", 100000, func(lex *jsonic.Lex, rule *jsonic.Rule) *jsonic.Token {
		pnt := lex.Cursor()
		src := lex.Src
		sI := pnt.SI
		srcLen := pnt.Len
		if sI >= srcLen {
			return nil
		}

		type kw struct {
			text string
			val  any
		}
		keywords := []kw{
			{"false", false},
			{"true", true},
			{"null", nil},
		}

		for _, k := range keywords {
			end := sI + len(k.text)
			if end > srcLen {
				continue
			}
			if src[sI:end] != k.text {
				continue
			}
			// Verify keyword boundary (not part of a longer identifier).
			if end < srcLen {
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

	// Trailing comma support: prepend close alternatives for pair and elem
	// rules so that ",}" and ",]" are accepted before the regular "," alt.
	if allowTrailingComma {
		CA := j.Token("#CA")
		CB := j.Token("#CB")
		CS := j.Token("#CS")

		j.Rule("pair", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{
				S: [][]jsonic.Tin{{CA}, {CB}},
				B: 1,
			})
		})

		j.Rule("elem", func(rs *jsonic.RuleSpec) {
			rs.PrependClose(&jsonic.AltSpec{
				S: [][]jsonic.Tin{{CA}, {CS}},
				B: 1,
			})
		})
	}

	// Add ZZ alt to val rule for empty/comment-only input.
	// Done programmatically to avoid Grammar() interfering with existing rules.
	j.Rule("val", func(rs *jsonic.RuleSpec) {
		rs.AddOpen(&jsonic.AltSpec{
			S: [][]jsonic.Tin{{ZZ}},
			G: "jsonc",
		})
	})
}

// ---- Options helpers ----

func optionsToMap(o *JsoncOptions) map[string]any {
	m := make(map[string]any)
	m["allowTrailingComma"] = boolOpt(o.AllowTrailingComma, false)
	m["disallowComments"] = boolOpt(o.DisallowComments, false)
	return m
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

// ---- Grammar helpers (shared pattern with ini/csv) ----

func mapToGrammarSpec(parsed map[string]any, ref map[jsonic.FuncRef]any) *jsonic.GrammarSpec {
	gs := &jsonic.GrammarSpec{
		Ref: ref,
	}

	ruleMap, _ := parsed["rule"].(map[string]any)
	if ruleMap == nil {
		return gs
	}

	gs.Rule = make(map[string]*jsonic.GrammarRuleSpec, len(ruleMap))
	for name, rDef := range ruleMap {
		rd, ok := rDef.(map[string]any)
		if !ok {
			continue
		}
		grs := &jsonic.GrammarRuleSpec{}
		if openDef, ok := rd["open"]; ok {
			grs.Open = convertAlts(openDef)
		}
		if closeDef, ok := rd["close"]; ok {
			grs.Close = convertAlts(closeDef)
		}
		gs.Rule[name] = grs
	}

	return gs
}

func convertAlts(def any) any {
	switch v := def.(type) {
	case []any:
		return convertAltList(v)
	case map[string]any:
		result := &jsonic.GrammarAltListSpec{}
		if alts, ok := v["alts"].([]any); ok {
			result.Alts = convertAltList(alts)
		}
		if inj, ok := v["inject"].(map[string]any); ok {
			result.Inject = &jsonic.GrammarInjectSpec{}
			if app, ok := inj["append"].(bool); ok {
				result.Inject.Append = app
			}
		}
		return result
	}
	return nil
}

func convertAltList(alts []any) []*jsonic.GrammarAltSpec {
	result := make([]*jsonic.GrammarAltSpec, 0, len(alts))
	for _, a := range alts {
		if am, ok := a.(map[string]any); ok {
			result = append(result, convertAlt(am))
		}
	}
	return result
}

func convertAlt(m map[string]any) *jsonic.GrammarAltSpec {
	ga := &jsonic.GrammarAltSpec{}

	if s, ok := m["s"]; ok {
		switch sv := s.(type) {
		case string:
			ga.S = sv
		case []any:
			strs := make([]string, len(sv))
			for i, v := range sv {
				strs[i], _ = v.(string)
			}
			ga.S = strs
		}
	}
	if b, ok := m["b"]; ok {
		ga.B = b
	}
	if p, ok := m["p"].(string); ok {
		ga.P = p
	}
	if r, ok := m["r"].(string); ok {
		ga.R = r
	}
	if a, ok := m["a"].(string); ok {
		ga.A = a
	}
	if c, ok := m["c"]; ok {
		ga.C = c
	}
	if e, ok := m["e"].(string); ok {
		ga.E = e
	}
	if g, ok := m["g"].(string); ok {
		ga.G = g
	}
	if u, ok := m["u"].(map[string]any); ok {
		ga.U = u
	}

	return ga
}
