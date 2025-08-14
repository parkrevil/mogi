import { InjectionToken } from '@nestjs/common';

export function makeServiceInjectionToken(name: string): InjectionToken {
  return `REDIS_SERVICE_${name.toUpperCase()}`;
}

export function makeClientInjectionToken(name: string): InjectionToken {
  return `REDIS_CLIENT_${name.toUpperCase()}`;
}