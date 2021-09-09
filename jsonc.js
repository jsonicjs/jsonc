"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Jsonc = void 0;
function Jsonc(jsonic, options) {
    jsonic.options({
        text: {
            lex: false,
        },
        number: {
            hex: false,
            oct: false,
            bin: false,
            sep: null,
        },
        string: {
            chars: '"',
            multiChars: '',
            allowUnknown: false,
            escape: {
                v: null,
            },
        },
        comment: {
            lex: true !== options.disallowComments,
        },
        map: {
            extend: false,
        },
        lex: {
            empty: false,
        },
        rule: {
            finish: false,
            include: 'json' + (options.allowTrailingComma ? ',comma' : ''),
        },
    });
}
exports.Jsonc = Jsonc;
//# sourceMappingURL=jsonc.js.map