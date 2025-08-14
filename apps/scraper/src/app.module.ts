import { Module } from '@nestjs/common';
import { ScrapingModule } from './scraping';
import { ConfigModule } from '@nestjs/config';
import { makeConfigModuleOptions } from '@mogi/bun-shared/configs';
import { mongoConfig } from '@mogi/bun-shared/configs';
import { MongooseModule } from '@nestjs/mongoose';
import { makeMongoModuleOptions } from '@mogi/bun-shared/providers/mongo';

@Module({
  imports: [
    ConfigModule.forRoot(
      makeConfigModuleOptions([mongoConfig])
    ),
    MongooseModule.forRootAsync(makeMongoModuleOptions()),
    ScrapingModule,
  ],
})
export class AppModule {}
