import { ClientSideCastlesLister } from "@find-castles/components/castles-lister/castles-lister.component";
import { Country, CountryCode, countries } from "@find-castles/lib/country";
import { getCastles } from "@find-castles/lib/db/get-castles";
import { MetadataProps } from "@find-castles/lib/metadata-props";
import { Box, Container, Typography } from "@mui/material";
import { ResolvingMetadata, Metadata } from "next";

const siteHost = process.env.SITE_HOST;

export async function generateMetadata(
  { params, searchParams }: MetadataProps,
  parent: ResolvingMetadata
): Promise<Metadata> {
  return {
    title: "Castles - Find Castles",
    description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
    keywords: ["castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
    applicationName: 'Find Castles',
    robots: { index: true, follow: true },
    authors: {
      name: 'Aurelio Buarque',
      url: 'https://buarki.com'
    },
    openGraph: {
      title: "Castles - Find Castles",
      description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
      url: `${siteHost}/castles`,
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
      title: "Castles - Find Castles",
      description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
      images: `${siteHost}/og.png`,
    }
  };
};


type CastlesPageProps = {
  searchParams: { [key: string]: string | string[] | undefined,
   },
};

// CSR+SSR
export default async function CastlesPage({ searchParams }: CastlesPageProps) {
  const countryCode = searchParams.country as CountryCode;
  const availableCoutries = countries.filter((country: Country) => country.trackingStatus === 'tracked');

  // const foundCastles = await getCastles({ contryCodes: ['uk'] });

  return (
    <Box sx={{ mt: 3, }}>
      <Container>
        <Typography variant="h1" color='primary.main'>Castles</Typography>
        <ClientSideCastlesLister countries={availableCoutries} currentCountry={countryCode}/>
      </Container>
    </Box>
  );
}

