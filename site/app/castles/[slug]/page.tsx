import { Box, Card, CardContent, CardMedia, Container, Grid, IconButton, Link, Tooltip, Typography } from "@mui/material";
import PhoneIcon from '@mui/icons-material/Phone';
import EmailIcon from '@mui/icons-material/Email';
import PetsIcon from '@mui/icons-material/Pets';
import LocalCafeIcon from '@mui/icons-material/LocalCafe';
import RestroomIcon from '@mui/icons-material/FamilyRestroom';
import StoreIcon from '@mui/icons-material/Store';
import PicnicIcon from '@mui/icons-material/Deck';
import LocalParkingIcon from '@mui/icons-material/LocalParking';
import MuseumIcon from '@mui/icons-material/Museum';
import AccessibleIcon from '@mui/icons-material/Accessible';
import { MetadataProps } from "@find-castles/lib/metadata-props";
import { ResolvingMetadata, Metadata } from "next";
import { getCastle } from "@find-castles/lib/db/getCastle";
import { CountryCode } from "@find-castles/lib/country";
import { decodeCastleURL } from "@find-castles/lib/encode-decore-url";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { Castle } from "@find-castles/lib/db/model";

export const dynamicParams = true;

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
      title: "Castles - Find Castles",
      description: "Explore the tracked and untracked European countries for castle data. Help us expand our data sources by contributing to the project.",
      images: `${siteHost}/og.png`,
    }
  };
};

interface Facilities {
  assistanceDogsAllowed: boolean;
  cafe: boolean;
  restrooms: boolean;
  giftshops: boolean;
  pinicArea: boolean;
  parking: boolean;
  exhibitions: boolean;
  wheelchairSupport: boolean;
}

const facilitiesIcons: { [key in keyof Facilities]: { icon: JSX.Element, label: string } } = {
  assistanceDogsAllowed: { icon: <PetsIcon />, label: "Assistance Dogs Allowed" },
  cafe: { icon: <LocalCafeIcon />, label: "Cafe" },
  restrooms: { icon: <RestroomIcon />, label: "Restrooms" },
  giftshops: { icon: <StoreIcon />, label: "Gift Shops" },
  pinicArea: { icon: <PicnicIcon />, label: "Picnic Area" },
  parking: { icon: <LocalParkingIcon />, label: "Parking" },
  exhibitions: { icon: <MuseumIcon />, label: "Exhibitions" },
  wheelchairSupport: { icon: <AccessibleIcon />, label: "Wheelchair Support" },
};

type CastlePageProps = {
  params: {
    slug: string;
  },
};

export default async function CastlePage({ params }: CastlePageProps) {
  const { countryCode, castleName } = decodeCastleURL(params.slug);

  const foundCastle = await getCastle({ name: castleName, country: countryCode as CountryCode });
  console.log({ foundCastle });
  console.log({ wh: JSON.stringify(foundCastle.visitingInfo) });

  const {
    name,
    pictureURL,
    contact,
    district,
    city,
    state,
    visitingInfo,
    sources,
    coordinates,
  } = foundCastle;

  return (
    <Container>
      <Box sx={{ my: 3, }}>
        <Typography variant="h4" align="center" gutterBottom>
          {toTitleCase(name)} Castle
        </Typography>

        <Card>
          <CardMedia
            component="img"
            height="400"
            image={pictureURL}
            alt={castleName}
          />
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Typography variant="subtitle1" gutterBottom>
                {
                  coordinates ? (
                    <Link title="See it on Google Maps" href={`https://www.google.com/maps/?q=${coordinates}`}>
                      {district}, {city}, {state}
                    </Link>
                  ) : <Typography>{district}, {city}, {state}</Typography>
                }
              </Typography>
              <Typography>{countryCode}</Typography>
            </Box>

            <Typography variant="h6" align="center" gutterBottom>
              Facilities
            </Typography>

            <Grid container spacing={2} justifyContent="center">
              {
                Object.values(visitingInfo?.facilities ?? {}).some((facility) => facility) ?
                Object.keys(visitingInfo?.facilities ?? {}).map((facility) => {
                  const key = facility as keyof Facilities;
                  return (visitingInfo?.facilities ?? {})[key] && (
                    <Grid item key={key} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                      <Tooltip title={facilitiesIcons[key]?.label} arrow>
                        <IconButton aria-label={facilitiesIcons[key]?.label}>
                          {facilitiesIcons[key]?.icon}
                        </IconButton>
                      </Tooltip>
                      <Typography variant="caption">{facilitiesIcons[key]?.label}</Typography>
                    </Grid>
                  );
                }) : <Typography align="center" sx={{ mt: 2 }}>No Info Available</Typography>
              }
            </Grid>

            <Box mt={2} display="flex" justifyContent="center">
              {contact?.phone && (
                <IconButton href={`tel:${contact.phone}`} aria-label="Phone">
                  <PhoneIcon />
                </IconButton>
              )}
              {contact?.email && (
                <IconButton href={`mailto:${contact.email}`} aria-label="Email">
                  <EmailIcon />
                </IconButton>
              )}
            </Box>
          </CardContent>
          {
            visitingInfo?.workingHours && (
              <Box sx={{ mt: 3 }}>
                <Typography variant="h6" align="center" gutterBottom>Working Hours</Typography>
                <Typography align="center" >
                  {visitingInfo?.workingHours ?? ''}
                </Typography>
              </Box>
            )
          }
        </Card>


        <Box sx={{ mt: 3 }}>
          <Typography variant="h4">Sources Of Information</Typography>
          <Typography>You may want to visit the following links to get more detailed data.</Typography>
          <ul>
            {
              sources.map((source: string, index: number) => (
                <li key={index}>
                  <a target="_blank" rel="noopener noreferrer" href={source}>{source}</a>
                </li>
              ))
            }
          </ul>
        </Box>
      </Box>
    </Container>
  );
}
