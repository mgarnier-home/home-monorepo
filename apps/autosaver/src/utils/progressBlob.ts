import { Blob } from 'fetch-blob';

export class ProgressBlob extends Blob {
  private progressCb?: (streamActual: number, streamTotal: number) => void;

  constructor(
    blobParts?: BlobPart[],
    options?: BlobPropertyBag,
    progressCb?: (streamActual: number, streamTotal: number) => void
  ) {
    super(blobParts, options);

    this.progressCb = progressCb;
  }

  stream(): ReadableStream<Uint8Array> {
    const progressCb = this.progressCb;
    const stream = super.stream();
    const streamTotal = this.size;

    let streamActual = 0;

    const stream2 = stream.pipeThrough(
      new TransformStream({
        transform(chunk, controller) {
          controller.enqueue(chunk);

          streamActual += chunk.length;

          if (progressCb) progressCb(streamActual, streamTotal);
        },
      })
    );

    return stream2;
  }
}
