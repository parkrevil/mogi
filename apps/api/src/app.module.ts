import { MongoConfig, RedisConfig, SharedConfig, makeConfigModuleOptions, mongoConfig, redisConfig } from '@mogi/bun-shared/configs';
import { isLocal } from '@mogi/bun-shared/helpers';
import { RedisModule } from '@mogi/bun-shared/modules/redis';
import { Module, OnApplicationShutdown } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { InjectConnection, MongooseModule } from '@nestjs/mongoose';
import type { Connection as MongooseConnection } from 'mongoose';
import { AppController } from './app.controller';
import { configs } from './core/configs';

@Module({
  imports: [
    ConfigModule.forRoot(
      makeConfigModuleOptions([...configs, mongoConfig, redisConfig])
    ),
    MongooseModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: async (configService: ConfigService) => {
        const config = configService.get<MongoConfig>(SharedConfig.Mongo);

        return {
          uri: config.uri,
          autoIndex: isLocal(),
        };
      },
      inject: [ConfigService],
    }),
    RedisModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: async (configService: ConfigService) => {
        const config = configService.get<RedisConfig>(SharedConfig.Redis);

        return {
          uri: config.uri,
        };
      },
      inject: [ConfigService],
    }),
  ],
  controllers: [AppController],
})
export class AppModule implements OnApplicationShutdown {
  constructor(@InjectConnection() private readonly mongooseConnection: MongooseConnection) { }

  async onApplicationShutdown() {
    console.log('Shutting down gracefully...');
    console.log('Closing Mongoose connection...');

    await this.mongooseConnection.close();

    console.log('Mongoose connection closed');
  }
}
