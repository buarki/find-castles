import { CountrySelector } from "@find-castles/components/country-selector/country-selector.component";
import { Country, countries } from "@find-castles/lib/country";
import { MetadataProps } from "@find-castles/lib/metadata-props";
import { Box, Container, Typography } from "@mui/material";
import { Metadata, ResolvingMetadata } from "next";

const siteHost = process.env.SITE_HOST;

const description = "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.";

export async function generateMetadata(
  { params, searchParams }: MetadataProps,
  parent: ResolvingMetadata
): Promise<Metadata> {
  return {
    title: "Find Castles",
    description: description,
    keywords: ["castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
    applicationName: 'Find Castles',
    robots: { index: true, follow: true },
    authors: {
      name: 'Aurelio Buarque',
      url: 'https://buarki.com'
    },
    openGraph: {
      title: "Find Castles",
      description: description,
      url: `${siteHost}`,
      type: "website",
      images: [
        {
          url: `${siteHost}/og.png`,
          width: 1200,
          height: 630,
          alt: "Find Castles",
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      site: "@buarki",
      title: "Find Castles",
      description: description,
      images: `${siteHost}/og.png`,
    }
  };
};

const trackedCountries = countries.filter((country: Country) => country.trackingStatus === 'tracked');

export default function Home() {
  return (
    <Box>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify(
            {
              "@context": "https://schema.org",
              "@graph": [
                {
                  "@type": "CollectionPage",
                  "@id": `${siteHost}/`,
                  "url": `${siteHost}/`,
                  "name": "Find Castles",
                  "isPartOf": {
                    "@id": `${siteHost}/#website`
                  },
                  "about": {
                    "@id": `${siteHost}/#/schema/person/3ug424jgiuf424`
                  },
                  "description": description,
                  "keywords": ["castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
                  "breadcrumb": {
                    "@id": `${siteHost}/#breadcrumb`
                  },
                  "inLanguage": "en-US"
                },
                {
                  "@type": "BreadcrumbList",
                  "@id": `${siteHost}/#breadcrumb`,
                  "itemListElement": [
                    {
                      "@type": "ListItem",
                      "position": 1,
                      "name": "Home",
                      "item": `${siteHost}/`
                    },
                    {
                      "@type": "ListItem",
                      "position": 2,
                      "name": "Castles",
                      "item": `${siteHost}/castles`
                    },
                    {
                      "@type": "ListItem",
                      "position": 3,
                      "name": "About",
                      "item": `${siteHost}/about`
                    },
                    {
                      "@type": "ListItem",
                      "position": 4,
                      "name": "Data Sources",
                      "item": `${siteHost}/data-sources`
                    },
                    {
                      "@type": "ListItem",
                      "position": 5,
                      "name": "Tech Stuff",
                      "item": `${siteHost}/tech-stuff`
                    }
                  ]
                },
                {
                  "@type": "WebSite",
                  "@id": `${siteHost}/#website`,
                  "url": `${siteHost}/`,
                  "name": "Find Castles",
                  "description": description,
                  "publisher": {
                    "@id": `${siteHost}/#/schema/person/3ug424jgiuf424`
                  },
                  "potentialAction": [
                    {
                      "@type": "SearchAction",
                      "target": {
                        "@type": "EntryPoint",
                        "urlTemplate": `${siteHost}/?s={search_term_string}`
                      },
                      "query-input": "required name=search_term_string"
                    }
                  ],
                  "inLanguage": "en-US"
                },
                {
                  "@type": [
                    "Person",
                    "Organization"
                  ],
                  "@id": `${siteHost}/#/schema/person/3ug424jgiuf424`,
                  "name": "Find Castles",
                  "image": {
                    "@type": "ImageObject",
                    "inLanguage": "en-US",
                    "@id": `${siteHost}/#/schema/person/image/`,
                    "url": `${siteHost}/og.png`,
                    "contentUrl": `${siteHost}/og.png`,
                    "width": 450,
                    "height": 450,
                    "caption": "Find Castles"
                  },
                  "logo": {
                    "@id": `${siteHost}/#/schema/person/image/`
                  }
                }
              ]
            }
          )
        }}
      />
      <Box
        sx={{
          background: 'secondary.main',
          minHeight: 'calc(100vh - 64px)',
          display: {
            xs: 'none',
            md: 'flex',
          },
          alignItems: 'center',
          justifyContent: 'center',
          backgroundImage: `url('/background.webp')`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          bgcolor: 'secondary.main',
          p: 2,
        }}>
        <Container>
          <Box
            sx={{
              width: '50vh',
              display: 'flex',
              flexDirection: 'column',
              gap: 3,
              justifyContent: 'left'
            }}>
            <Typography variant="h1" color='primary.main'>Choose A Country</Typography>
            <Typography>Embark on a journey with Castle Explorer and uncover the secrets of these timeless treasures. Select the country to see the castles and start your adventure today!</Typography>

            <CountrySelector countries={trackedCountries} />
          </Box>
        </Container>
      </Box>
      <Box
        sx={{
          display: {
            xs: 'flex',
            md: 'none',
          },
          background: 'secondary.main',
          minHeight: 'calc(100vh - 64px)',
          alignItems: 'center',
          justifyContent: 'center',
          backgroundImage: `url('/background-mobile.webp')`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          bgcolor: 'secondary.main',
        }}>
        <Container>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, }}>
            <Typography variant="h1" color='primary.main'>Choose A Country</Typography>
            <Typography>Embark on a journey with Castle Explorer and uncover the secrets of these timeless treasures. Select the country to see the castles and start your adventure today!</Typography>

            <CountrySelector countries={trackedCountries} />
          </Box>
        </Container>
      </Box>
      <Container>
        <Box>
          <Typography></Typography>
        </Box>
      </Container>
    </Box>
  );
}
