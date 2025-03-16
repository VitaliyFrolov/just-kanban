// Removes all empty values from object by provided fields values
export function cleanObject<T extends Record<string, unknown>>(
  obj: T,
  excludedValues: unknown[] = [undefined, null, ''],
): Partial<T> {
  const cleanedObject: Partial<T> = {};

  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const value = obj[key];
      if (excludedValues.includes(value)) {
        cleanedObject[key] = value;
      }
    }
  }

  return cleanedObject;
}
