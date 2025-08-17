import { registerAs } from '@nestjs/config';
import { Config } from './enums';
import { HttpServerConfig } from './interfaces';

export default registerAs<HttpServerConfig>(Config.HttpServer, () => ({
  listening: {
    host: process.env.API_HTTP_SERVER_LISTENING_HOST,
    port: parseInt(process.env.API_HTTP_SERVER_LISTENING_PORT, 10),
  },
}));
