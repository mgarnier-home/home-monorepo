import { environment } from '../../environments/environment';

export class ApiError extends Error {
  constructor(message: any) {
    super(message);
    this.name = 'ApiError';
  }
}

export abstract class ApiService {
  constructor() {}

  protected getApiUrl(): string {
    return environment.apiUrl;
  }

  private async makeRequest<T>(url: string | URL | globalThis.Request, init?: RequestInit): Promise<T> {
    try {
      const response = await fetch(url, init);

      if (!response.ok) {
        throw response.body?.toString();
      }

      const data = await response.json();

      return data;
    } catch (error) {
      console.error(`New Api Error: ${error}`);

      throw new ApiError(error);
    }
  }

  protected async get<T>(url: string): Promise<T> {
    return this.makeRequest(url);
  }

  protected async postJson<T>(url: string, body: any): Promise<T> {
    return this.makeRequest(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });
  }

  protected async delete<T>(url: string): Promise<T> {
    return this.makeRequest(url, {
      method: 'DELETE',
    });
  }
}
