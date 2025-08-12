import { ConfigModule, ConfigService } from '@nestjs/config';
import { SharedConfig, IMongoConfig } from '@mogi/bun-shared/configs';
import { MongooseModuleAsyncOptions } from '@nestjs/mongoose';
import { isLocal } from '../../helpers';

export function makeMongoModuleOptions(): MongooseModuleAsyncOptions {
  return {
    imports: [ConfigModule],
    useFactory: async (configService: ConfigService) => {
      const mongoConfig = configService.get<IMongoConfig>(SharedConfig.MONGO);

      return {
        uri: mongoConfig!.uri,
        autoIndex: isLocal()
      };
    },
    inject: [ConfigService],
  };
}
