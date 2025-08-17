import { registerAs } from '@nestjs/config';
import { SharedConfig } from './enums';
import { RedisConfig } from './interfaces';

export default registerAs<RedisConfig>(SharedConfig.Redis, () => ({
  uri: process.env.REDIS_URI!,
}));
