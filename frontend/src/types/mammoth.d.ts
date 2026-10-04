declare module 'mammoth' {
  export interface ConvertResult {
    value: string;
    messages: Array<{
      type: string;
      message: string;
    }>;
  }

  export interface ConvertOptions {
    arrayBuffer?: ArrayBuffer;
    buffer?: any;
    path?: string;
    styleMap?: string | string[];
    includeDefaultStyleMap?: boolean;
    convertImage?: any;
  }

  export const images: {
    inline: (fn: (element: any) => Promise<{ src: string }>) => any;
  };

  export function convertToHtml(
    input: { arrayBuffer: ArrayBuffer } | { buffer: any } | { path: string },
    options?: ConvertOptions
  ): Promise<ConvertResult>;

  export function extractRawText(
    input: { arrayBuffer: ArrayBuffer } | { buffer: any } | { path: string },
    options?: ConvertOptions
  ): Promise<ConvertResult>;
}
