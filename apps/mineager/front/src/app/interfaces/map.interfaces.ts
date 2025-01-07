export interface Map {
  name: string;
  version: string;
  description: string;
}

export interface CreateMapRequest {
  name: string;
  version: string;
  description: string;
  file: string;
}
