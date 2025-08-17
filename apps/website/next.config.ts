import dotenv from 'dotenv';
import type { NextConfig } from "next";

dotenv.config({ path: `../../.env.${process.env.NODE_ENV}` })

const nextConfig: NextConfig = {
  env: {
    API_SERVER_URL: process.env.WEBSITE_API_SERVER_URL,
  },
  image: {
    formats: ['image/avif', 'image/webp'],
  }
};

export default nextConfig;
