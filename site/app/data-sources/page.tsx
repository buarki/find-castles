import { Country, countries } from "@find-castles/lib/country";
import { MetadataProps } from "@find-castles/lib/metadata-props";
import { Box, Container, List, ListItem, Typography, Divider, Link } from "@mui/material";
import { Metadata, ResolvingMetadata } from "next";

const siteHost = process.env.SITE_HOST;

const trackedCountries = countries.filter((country: Country) => country.trackingStatus === 'tracked');
const untrackedCountries = countries.filter((country: Country) => country.trackingStatus !== 'tracked');
const trackedCountriesPercentage = (trackedCountries.length / countries.length) * 100;
const text = `So far we have tracked ${trackedCountriesPercentage.toFixed(2)}% of European countries. The full list is available below.`;

export async function generateMetadata(
  { params, searchParams }: MetadataProps,
  parent: ResolvingMetadata
): Promise<Metadata> {
  return {
    title: "Data Sources - Find Castles",
    description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
    keywords: ["castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
    applicationName: 'Find Castles',
    robots: { index: true, follow: true },
    authors: {
      name: 'Aurelio Buarque',
      url: 'https://buarki.com'
    },
    openGraph: {
      title: "Data Sources - Find Castles",
      description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
      url: `${siteHost}/data-sources`,
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
      title: "Data Sources - Find Castles",
      description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
      images: `${siteHost}/og.png`,
    }
  };
};

export default function DataSourcesPage() {
  return (
    <Box sx={{ pt: 8, pb: 8, bgcolor: 'secondary.main' }}>
      <Container>
        <header>
          <Typography variant="h1" gutterBottom>Tracked and Untracked Countries</Typography>
          <Typography variant="h5" gutterBottom>{text}</Typography>
          <Typography variant="h5">You can help us to increase this number by creating a web scraper.</Typography>
        </header>

        <Divider sx={{ my: 4 }} />

        <section>
          <Typography variant="h2" gutterBottom>Tracked Countries</Typography>
          <List sx={{ pl: 2 }}>
            {trackedCountries.map((country: Country) => (
              <ListItem key={country.code} sx={{ pl: 0, py: 0.5 }}>
                <Link underline="always" href={`/castles?country=${country.code}`} sx={{ color: 'inherit' }}>
                  <Typography variant="body1">{country.name}</Typography>
                </Link>
              </ListItem>
            ))}
          </List>
        </section>

        <Divider sx={{ my: 4 }} />

        <section>
          <Typography variant="h2" gutterBottom>Untracked Countries</Typography>
          <List sx={{ pl: 2 }}>
            {untrackedCountries.map((country: Country) => (
              <ListItem key={country.code} sx={{ pl: 0, py: 0.5 }}>
                <Typography variant="body1">{country.name}</Typography>
              </ListItem>
            ))}
          </List>
        </section>
      </Container>
    </Box>
  );
}


