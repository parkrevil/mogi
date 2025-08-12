import { registerAs } from '@nestjs/config';
import { Config } from './enums';
import { IHttpServerConfig } from './interfaces';

export default registerAs<IHttpServerConfig>(Config.HTTP_SERVER, () => ({
  port: parseInt(process.env.API_HTTP_SERVER_PORT, 10),
  host: process.env.API_HTTP_SERVER_HOST,
}));
