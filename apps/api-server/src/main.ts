import { NestFactory } from '@nestjs/core';
import type { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import { Config, IHttpServerConfig } from './core/configs';

async function bootstrap() {
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  const configService = app.get(ConfigService);
  const httpServerConfig = configService.get<IHttpServerConfig>(Config.HTTP_SERVER);

  await app.listen(httpServerConfig.port, httpServerConfig.host);

  console.log(`Server is running on ${httpServerConfig.host}:${httpServerConfig.port}`);
}

bootstrap();
