export const SERVER_ROUTES = {
  CONF: '/conf',
  PING_HOST: '/ping-host',
  MAKE_REQUEST: '/make-request',
  STATUS_CHECKS: '/status-checks',
} as const;

export type ServerRouteKeys = keyof typeof SERVER_ROUTES;
export type ServerRoute = (typeof SERVER_ROUTES)[ServerRouteKeys];

export const SERVER_ROUTES_METHODS: Record<ServerRoute, 'GET' | 'POST'> = {
  [SERVER_ROUTES.CONF]: 'GET',
  [SERVER_ROUTES.PING_HOST]: 'POST',
  [SERVER_ROUTES.MAKE_REQUEST]: 'POST',
  [SERVER_ROUTES.STATUS_CHECKS]: 'POST',
} as const;
