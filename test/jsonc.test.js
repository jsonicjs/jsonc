"use strict";
/* Copyright (c) 2021 Richard Rodger and other contributors, MIT License */
Object.defineProperty(exports, "__esModule", { value: true });
const jsonic_1 = require("jsonic");
const jsonc_1 = require("../jsonc");
describe('jsonc', () => {
    test('happy', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        expect(j('{"a":1}')).toEqual({ a: 1 });
    });
    test('comments', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        expect(j('// this is a comment')).toEqual(undefined);
        expect(j('// this is a comment\n')).toEqual(undefined);
        expect(j('/* this is a comment*/')).toEqual(undefined);
        expect(j('/* this is a \r\ncomment*/')).toEqual(undefined);
        expect(j('/* this is a \ncomment*/')).toEqual(undefined);
        expect(() => j('/* this is a')).toThrow('unterminated_comment');
        expect(() => j('/* this is a \ncomment')).toThrow('unterminated_comment');
        expect(() => j('/ ttt')).toThrow('unexpected');
    });
    test('strings', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        // console.dir(j.options)
        // console.dir(j.internal().config)
        expect(j('"test"')).toEqual('test');
        expect(j('"\\""')).toEqual('"');
        expect(j('"\\/"')).toEqual('/');
        expect(j('"\\b"')).toEqual('\b');
        expect(j('"\\f"')).toEqual('\f');
        expect(j('"\\n"')).toEqual('\n');
        expect(j('"\\r"')).toEqual('\r');
        expect(j('"\\t"')).toEqual('\t');
        expect(j('"\u88ff"')).toEqual('\u88ff');
        expect(j('"​\u2028"')).toEqual('​\u2028');
        expect(() => j('"\\v"')).toThrow('unexpected');
        expect(() => j('"test')).toThrow('unterminated_string');
        expect(() => j('"test\n"')).toThrow('unprintable');
        expect(() => j('"\t"')).toThrow('unprintable');
        expect(() => j('"\t "')).toThrow('unprintable');
        expect(() => j('"\0 "')).toThrow('unprintable');
    });
    test('numbers', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        expect(j('0')).toEqual(0);
        expect(j('0.1')).toEqual(0.1);
        expect(j('-0.1')).toEqual(-0.1);
        expect(j('-1')).toEqual(-1);
        expect(j('1')).toEqual(1);
        expect(j('123456789')).toEqual(123456789);
        expect(j('10')).toEqual(10);
        expect(j('90')).toEqual(90);
        expect(j('90E+123')).toEqual(90E+123);
        expect(j('90e+123')).toEqual(90e+123);
        expect(j('90e-123')).toEqual(90e-123);
        expect(j('90E-123')).toEqual(90E-123);
        expect(j('90E123')).toEqual(90E123);
        expect(j('90e123')).toEqual(90e123);
        expect(j('01')).toEqual(1);
        expect(j('-01')).toEqual(-1);
        expect(() => j('-')).toThrow('unexpected');
        expect(() => j('.0')).toThrow('unexpected');
    });
    test('keywords', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        expect(j('true')).toEqual(true);
        expect(j('false')).toEqual(false);
        expect(j('null')).toEqual(null);
        expect(() => j('nulllll')).toThrow('unexpected');
        expect(() => j('True')).toThrow('unexpected');
        expect(() => j('foo-bar')).toThrow('unexpected');
        expect(() => j('foo bar')).toThrow('unexpected');
        expect(j('false//hello')).toEqual(false);
    });
    test('trivia', () => {
        let j = jsonic_1.Jsonic.make().use(jsonc_1.Jsonc);
        expect(j(' ')).toEqual(undefined);
        expect(j('  \t  ')).toEqual(undefined);
        expect(j('  \t  \n  \t  ')).toEqual(undefined);
        expect(j('\r\n')).toEqual(undefined);
        expect(j('\r')).toEqual(undefined);
        expect(j('\n')).toEqual(undefined);
        expect(j('\n\r')).toEqual(undefined);
        expect(j('\n   \n')).toEqual(undefined);
    });
});
//# sourceMappingURL=jsonc.test.js.map