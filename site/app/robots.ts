import { MetadataRoute } from 'next'
 
export default function robots(): MetadataRoute.Robots {
  const appHost = process.env.SITE_HOST;
  
  return {
    rules: [
      {
        userAgent: '*',
        disallow: " ",
      },
    ],
    sitemap: `${appHost}/sitemap.xml`,
  }
}
