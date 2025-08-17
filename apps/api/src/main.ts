import { bootstrap } from '@mogi/bun';
import { ConfigService } from '@nestjs/config';
import { NestFactory } from '@nestjs/core';
import type { NestExpressApplication } from '@nestjs/platform-express';
import { AppModule } from './app.module';
import { Config, HttpServerConfig } from './core/configs';

bootstrap(async () => {
  const app = await NestFactory.create<NestExpressApplication>(AppModule);
  const configService = app.get(ConfigService);
  const { listening } = configService.get<HttpServerConfig>(Config.HttpServer);

  await app.listen(listening.port, listening.host);

  console.log(`Server is running on ${listening.host}:${listening.port}`);

  return app;
});
