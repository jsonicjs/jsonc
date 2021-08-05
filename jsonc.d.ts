import { Jsonic } from 'jsonic';
declare type JsoncOptions = {
    allowTrailingComma?: boolean;
    disallowComments?: boolean;
};
declare function Jsonc(jsonic: Jsonic, options: JsoncOptions): void;
export { Jsonc, };
export type { JsoncOptions, };
