import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { MongooseModule } from '@nestjs/mongoose';
import { AppController } from './app.controller';
import { configs } from './core/configs';
import { mongoConfig, makeConfigModuleOptions } from '@mogi/bun-shared/configs';
import { makeMongoModuleOptions } from '@mogi/bun-shared/providers/mongo';

@Module({
  imports: [
    ConfigModule.forRoot(
      makeConfigModuleOptions([...configs, mongoConfig])
    ),
    MongooseModule.forRootAsync(makeMongoModuleOptions()),
  ],
  controllers: [AppController],
})
export class AppModule {}
