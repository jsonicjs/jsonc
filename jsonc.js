"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Jsonc = void 0;
function Jsonc(jsonic, options) {
    // TODO: merge these calls
    jsonic.options({
        text: { lex: false },
        number: { hex: false, oct: false, bin: false, sep: null },
        string: { chars: '"', multiChars: '', allowUnknown: false, escape: { v: null } },
        comment: { lex: false },
        map: { extend: false },
        lex: { empty: false },
        rule: {
            finish: false,
            include: 'json' + (options.allowTrailingComma ? ',comma' : '')
        },
    });
    jsonic.options({
        comment: {
            lex: true && !options.disallowComments,
            marker: [
                { line: true, start: '//', lex: true },
                { line: false, start: '/' + '*', end: '*' + '/', lex: true },
            ],
        }
    });
}
exports.Jsonc = Jsonc;
//# sourceMappingURL=jsonc.js.map