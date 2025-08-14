import { bootstrap } from '@mogi/bun-shared';
import { ConfigService } from '@nestjs/config';
import { NestFactory } from '@nestjs/core';
import type { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { Config, HttpServerConfig } from './core/configs';

bootstrap(async () => {
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  const configService = app.get(ConfigService);
  const httpServerConfig = configService.get<HttpServerConfig>(Config.HttpServer);

  await app.listen(httpServerConfig.port, httpServerConfig.host);

  console.log(`Server is running on ${httpServerConfig.host}:${httpServerConfig.port}`);

  return app;
});
