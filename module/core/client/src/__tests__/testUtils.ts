export function requireNonNull<T>(value: T | null, label: string): T {
  if (value === null) {
    throw new Error(`expected ${label}`);
  }
  return value;
}
