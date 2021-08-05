"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Jsonc = void 0;
function Jsonc(jsonic, options) {
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