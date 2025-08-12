import { Environment } from './enums';

export function isLocal() {
  return process.env.NODE_ENV === Environment.Local;
}

export function isProduction() {
  return process.env.NODE_ENV === Environment.Production;
}
