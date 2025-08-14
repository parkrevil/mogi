import { Module } from '@nestjs/common';
import { ScrapingService } from './scraping.service';

@Module({
  providers: [ScrapingService],
})
export class ScrapingModule {}
