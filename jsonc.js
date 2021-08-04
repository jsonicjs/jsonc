"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Jsonc = void 0;
function Jsonc(jsonic) {
    jsonic.options({
        // TODO: replace with a Jsonic.options('json') preset
        text: { lex: false },
        number: { hex: false, oct: false, bin: false, sep: null },
        string: { chars: '"', multiChars: '', allowUnknown: false, escape: { v: null } },
        comment: { lex: false },
        map: { extend: false },
        rule: { include: 'json' },
    });
    jsonic.options({
        comment: {
            lex: true,
            marker: [
                { line: true, start: '//', lex: true },
                { line: false, start: '/' + '*', end: '*' + '/', lex: true },
            ]
        }
    });
}
exports.Jsonc = Jsonc;
//# sourceMappingURL=jsonc.js.map