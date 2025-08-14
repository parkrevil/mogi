import { registerAs } from '@nestjs/config';
import { Config } from './enums';
import { HttpServerConfig } from './interfaces';

export default registerAs<HttpServerConfig>(Config.HttpServer, () => ({
  port: parseInt(process.env.API_HTTP_SERVER_PORT, 10),
  host: process.env.API_HTTP_SERVER_HOST,
}));
