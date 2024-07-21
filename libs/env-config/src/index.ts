export function getEnvVariable<T extends string | number | boolean>(
  variableName: string,
  required: boolean = false,
  defaultValue?: T
): T {
  let value: any = process.env[variableName];

  if (required === true) {
    if (typeof value === 'undefined') {
      throw new Error(`Environment variable ${variableName} is required, but was not provided.`);
    }

    if (typeof value === 'string' && value === '') {
      throw new Error(`Environment variable ${variableName} is required, but was provided with an empty value.`);
    }
  }

  if (typeof value === 'undefined') {
    value = defaultValue;
  }

  if (typeof defaultValue === 'number') {
    value = Number(value as string);

    if (isNaN(value)) {
      throw new Error(`Environment variable ${variableName} is not a number.`);
    }
  }

  if (typeof defaultValue === 'boolean') {
    if (value !== 'true' && value !== 'false') {
      throw new Error(`Environment variable ${variableName} is neither 'true' nor 'false'.`);
    }

    value = value === 'true';
  }

  return value as T;
}
