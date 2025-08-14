import { registerAs } from '@nestjs/config';
import { SharedConfig } from './enums';
import { MongoConfig } from './interfaces';

export default registerAs<MongoConfig>(SharedConfig.Mongo, () => ({
  uri: process.env.MONGO_URI!,
}));
