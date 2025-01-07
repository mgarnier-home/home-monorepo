export interface Server {
  name: string;
  url: string;
  version: string;
  map: string;
  status: string;
}

export interface CreateServerRequest {
  name: string;
  version: string;
  newMap: boolean;
  mapName: string;
  memory: string;
}

export interface DeleteServerRequest {
  full: boolean;
}
