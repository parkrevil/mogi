import { Inject } from '@nestjs/common';
import { makeClientInjectionToken, makeServiceInjectionToken } from './helpers';

export const InjectRedisService = (name: string) => Inject(makeServiceInjectionToken(name));
export const InjectRedisClient = (name: string) => Inject(makeClientInjectionToken(name));
