import { Injectable, OnModuleInit, OnApplicationShutdown } from '@nestjs/common';
import { RedisClient } from 'bun';

@Injectable()
export class RedisService implements OnModuleInit, OnApplicationShutdown {
  constructor(private readonly client: RedisClient) { }

  async onModuleInit() {
    await this.connect();
  }

  async onApplicationShutdown() {
    this.close();
  }

  private async connect() {
    try {
      await this.client.connect();
    } catch (e) {
      console.error('Failed to connect to Redis:', e);

      throw e;
    }
  }

  private close() {
    console.log('Closing Redis connection...');
    this.client.close();
    console.log('Redis connection closed');
  }
}
