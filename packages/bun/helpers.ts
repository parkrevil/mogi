import { INestApplication } from '@nestjs/common';
import { Environment } from './enums';

export async function bootstrap(starter: () => Promise<INestApplication>) {
  const exceptionEvents = ['uncaughtException', 'unhandledRejection'];

  await starter().then(app => {
    app.enableShutdownHooks();
  });

  for (const event of exceptionEvents) {
    process.on(event, async (e) => {
      console.error(`${event} occurred:`, e);
    });
  }
}

export function isLocal() {
  return process.env.NODE_ENV === Environment.Local;
}

export function isProduction() {
  return process.env.NODE_ENV === Environment.Production;
}
