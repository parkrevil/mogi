import { ConfigFactory } from '@nestjs/config';

export function makeConfigModuleOptions(configs: Array<ConfigFactory | Promise<ConfigFactory>>) {
  return {
    isGlobal: true,
    envFilePath: [`../../.env.${process.env.NODE_ENV}`],
    load: configs,
  };
}
