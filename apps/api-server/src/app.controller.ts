import { Controller, Get } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Controller()
export class AppController {
  constructor(private readonly configService: ConfigService) { }

  @Get()
  getHello(): string {
    return 'Hello World from API Server!';
  }

  @Get('health')
  getHealth() {
    return {
      status: 'ok',
      service: 'api-server',
      environment: this.configService.get('NODE_ENV', 'development'),
      timestamp: new Date().toISOString(),
    };
  }

  @Get('config')
  getConfig() {
    return {
      port: this.configService.get('PORT', 3000),
      environment: this.configService.get('NODE_ENV', 'development'),
      database: {
        host: this.configService.get('DB_HOST', 'localhost'),
        port: this.configService.get('DB_PORT', 27017),
        database: this.configService.get('DB_DATABASE', 'mogi'),
      },
      redis: {
        host: this.configService.get('REDIS_HOST', 'localhost'),
        port: this.configService.get('REDIS_PORT', 6379),
        db: this.configService.get('REDIS_DB', 0),
      },
    };
  }
}
