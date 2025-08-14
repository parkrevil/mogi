import { DynamicModule, Module } from '@nestjs/common';
import { RedisClient } from 'bun';
import { makeClientInjectionToken, makeServiceInjectionToken } from './helpers';
import { RedisModuleAsyncOptions } from './interfaces';
import { RedisService } from './redis.service';

@Module({})
export class RedisModule {
  private static readonly DEFAULT_NAME = 'default';
  private static services = new Map<string, RedisService>();
  private static clients = new Map<string, RedisClient>();

  static forRootAsync(options: RedisModuleAsyncOptions): DynamicModule {
    const name = (options.name ? String(options.name) : RedisModule.DEFAULT_NAME).toUpperCase();
    const serviceToken = makeServiceInjectionToken(name);
    const clientToken = makeClientInjectionToken(name);

    return {
      module: RedisModule,
      imports: options.imports || [],
      providers: [
        {
          provide: clientToken,
          useFactory: async (...params: any[]) => {
            const { uri, ...redisOptions } = await options.useFactory(...params);

            if (RedisModule.clients.has(name)) {
              return RedisModule.clients.get(name);
            }

            const client = new RedisClient(uri, redisOptions);

            client.onconnect = () => {
              console.log('Connected to Redis server');
            };
            client.onclose = e => {
              console.error('Disconnected from Redis server:', e);
            };

            RedisModule.clients.set(name, client);

            return client;
          },
          inject: options.inject || [],
        },
        {
          provide: serviceToken,
          useFactory: async (client: RedisClient, ...params: any[]) => {
            const key = (name || RedisModule.DEFAULT_NAME).toString();

            if (RedisModule.services.has(key)) {
              return RedisModule.services.get(key);
            }

            const service = new RedisService(client);

            RedisModule.services.set(key, service);

            return service;
          },
          inject: [clientToken].concat(options.inject || []),
        },
      ],
      exports: [serviceToken, clientToken],
    };
  }
}
