import { registerAs } from '@nestjs/config';
import { SharedConfig } from './enums';
import { IMongoConfig } from './interfaces';

export default registerAs<IMongoConfig>(SharedConfig.MONGO, () => ({
  uri: process.env.MONGO_URI!,
}));
