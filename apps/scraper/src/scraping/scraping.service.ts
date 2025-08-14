import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ScrapingService {
  constructor(private readonly configService: ConfigService) {}
}
