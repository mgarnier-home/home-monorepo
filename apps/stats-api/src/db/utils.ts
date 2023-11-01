export namespace DbUtils {
  export function groupBy<T>(list: T[], getKey: (item: T) => string): Record<string, T[]> {
    const map: Record<string, T[]> = {};

    for (const item of list) {
      const key = getKey(item);
      const collection: T[] = map[key] || (map[key] = []);
      collection.push(item);
    }

    return map;
  }
}
