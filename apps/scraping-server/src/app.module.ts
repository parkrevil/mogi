import { Module } from '@nestjs/common';
import { ScrapingModule } from './scraping';

@Module({
  imports: [ScrapingModule],
})
export class AppModule {}
