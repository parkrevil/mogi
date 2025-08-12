import { registerAs } from '@nestjs/config';
import { Config } from './enums';

export default registerAs(Config.HTTP_SERVER, () => ({
  port: parseInt(process.env.HTTP_SERVER_PORT, 10),
  host: process.env.HTTP_SERVER_HOST,
}));
