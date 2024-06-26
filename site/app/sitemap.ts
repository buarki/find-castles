import { getCastles } from "@find-castles/lib/db/get-castles";
import { Castle } from "@find-castles/lib/db/model";
import { MetadataRoute } from "next";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const appHost = process.env.SITE_HOST;;
  
  const pages = ['', 'about', 'castles', 'data-sources', 'tech-stuff'].map(
    (route) => ({
      url: `${appHost}/${route}`,
      lastModified: new Date().toISOString().split('T')[0],
    })
  );
  const castlePages = (await getCastles({ contryCodes: ['ie', 'pt', 'uk', 'sk'] })).map((foundCastle: Castle) => ({
    url: `${appHost}/castles/${foundCastle.webName}`,
    lastModified: new Date().toISOString().split('T')[0],
  }));
  return [
    ...pages,
    ...castlePages,
  ];
}
