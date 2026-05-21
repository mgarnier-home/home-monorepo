export class SimpleCache<T> {
  private data: Map<string, { timestamp: number; value: T }> = new Map();
  private ttl: number; // Time-to-live in milliseconds

  constructor(ttlSeconds: number = 15) {
    this.ttl = ttlSeconds * 1000; // Convert to milliseconds
  }

  set(key: string, value: T): void {
    const timestamp = Date.now();
    this.data.set(key, { timestamp, value });
  }

  get(key: string): T | null {
    const cached = this.data.get(key);

    if (!cached) return null;

    const now = Date.now();
    if (now - cached.timestamp > this.ttl) {
      this.data.delete(key);
      return null;
    }

    return cached.value;
  }

  has(key: string): boolean {
    return this.data.has(key);
  }

  invalidate(key: string): void {
    this.data.delete(key);
  }
}
