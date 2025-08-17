import { makeConfigModuleOptions, mongoConfig } from '@mogi/bun/configs';
import { makeMongoModuleOptions } from '@mogi/bun/providers/mongo';
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { MongooseModule } from '@nestjs/mongoose';
import { ScrapingModule } from './scraping';

@Module({
  imports: [
    ConfigModule.forRoot(
      makeConfigModuleOptions([mongoConfig])
    ),
    MongooseModule.forRootAsync(makeMongoModuleOptions()),
    ScrapingModule,
  ],
})
export class AppModule { }
