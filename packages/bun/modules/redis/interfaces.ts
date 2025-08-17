
export interface RedisModuleAsyncOptions {
  name?: string | symbol;
  imports?: any[];
  useFactory: (...args: any[]) => RedisOptions | Promise<RedisOptions>;
  inject?: any[];
}

export interface RedisOptions extends Bun.RedisOptions {
  uri: string;
}
