import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppController } from './app.controller';
import { configs } from './core/configs';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: [`../../.env.${process.env.NODE_ENV}`],
      load: [...configs],
    }),
  ],
  controllers: [AppController],
})
export class AppModule {}
