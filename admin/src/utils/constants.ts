export const APP_NAME = 'Mynute Admin';
export const API_BASE_URL = '/api';
export const TOKEN_KEY = 'admin_token';

export const ROUTES = {
    HOME: '/',
    DASHBOARD: '/',
    USERS: '/users',
    COMPANIES: '/companies',
    SETTINGS: '/settings',
} as const;

export const HTTP_STATUS = {
    OK: 200,
    CREATED: 201,
    BAD_REQUEST: 400,
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    SERVER_ERROR: 500,
} as const;
