import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  allowedDevOrigins: ['172.19.0.1'],
  async rewrites() {
    return [
      {
        source: '/static/:path*',
        destination: 'http://localhost:8082/static/:path*',
      },
    ];
  },
};

export default nextConfig;
