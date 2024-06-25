import { Box, Card, CardContent, CardMedia, Container, Grid, IconButton, Link, Tooltip, Typography, List, ListItem, ListItemText, ListItemIcon } from "@mui/material";
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
import InfoIcon from '@mui/icons-material/Info';
import { MetadataProps } from "@find-castles/lib/metadata-props";
import { ResolvingMetadata, Metadata } from "next";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { getCastles } from "@find-castles/lib/db/get-castles";
import { getCastle } from "@find-castles/lib/db/get-castle";
import { notFound } from "next/navigation";

export const dynamicParams = true;

const siteHost = process.env.SITE_HOST!;

export async function generateMetadata(
  { params, searchParams }: MetadataProps,
  parent: ResolvingMetadata
): Promise<Metadata> {
  const castleWebName = params.slug;
  const foundCastle = await getCastle(castleWebName);
  if (!foundCastle) {
    notFound();
  }

  return {
    metadataBase: new URL('https://find-castles.vercel.app'),
    title: toTitleCase(foundCastle.name),
    description: `Discover ${foundCastle.name} castle on Find Castles`,
    keywords: [foundCastle.name, foundCastle.country, "castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
    applicationName: 'Find Castles',
    robots: { index: true, follow: true },
    authors: {
      name: 'Aurelio Buarque',
      url: 'https://buarki.com'
    },
    openGraph: {
      title: toTitleCase(foundCastle.name),
      description: `Discover ${foundCastle.name} castle on Find Castles`,
      url: `${siteHost}/${foundCastle.webName}`,
      type: "website",
      images: [
        {
          url: foundCastle.pictureURL,
          width: 1200,
          height: 630,
          alt: "Find Castles",
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      site: "@buarki",
      title: toTitleCase(foundCastle.name),
      description: `Discover ${foundCastle.name} castle on Find Castles`,
      images: foundCastle.pictureURL,
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

export async function generateStaticParams() {
  return (await getCastles({contryCodes: ['ie', 'pt', 'uk', 'sk']}))
          .map((foundCastle) => ({ slug: foundCastle.webName, }));
}

export default async function CastlePage({ params }: CastlePageProps) {
  const webName = params.slug;
  const foundCastle = await getCastle(webName);
  if (!foundCastle) {
    notFound();
  }

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
    country,
    propertyCondition,
  } = foundCastle;

  return (
    <Container>
      <Box sx={{ my: 3 }}>
        <Typography variant="h4" align="center" gutterBottom>
          {toTitleCase(name)} Castle
        </Typography>

        <Card>
          <CardMedia
            component="img"
            height="400"
            image={pictureURL}
            alt={name}
          />
          <CardContent>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Typography variant="subtitle1" gutterBottom>
                {
                  coordinates ? (
                    <Link title="See it on Google Maps" target="_blank" href={`https://www.google.com/maps/?q=${coordinates}`}>
                      {district}, {city}, {state}
                    </Link>
                  ) : <Typography>{district}, {city}, {state}</Typography>
                }
              </Typography>
              <Typography>{country}</Typography>
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

            <Box mt={3}>
              <Typography variant="h6" align="center" gutterBottom>
                Contact
              </Typography>
              <Box sx={{ display: 'flex', justifyContent: 'center', gap: 2 }}>
                {contact?.phone ? (
                  <IconButton href={`tel:${contact.phone}`} aria-label="Phone">
                    <Tooltip title="Phone" arrow>
                      <PhoneIcon />
                    </Tooltip>
                  </IconButton>
                ) : <Typography>No Phone Info</Typography>}
                {contact?.email ? (
                  <IconButton href={`mailto:${contact.email}`} aria-label="Email">
                    <Tooltip title="Email" arrow>
                      <EmailIcon />
                    </Tooltip>
                  </IconButton>
                ) : <Typography>No Email Info</Typography>}
              </Box>
            </Box>

            {visitingInfo?.workingHours && (
              <Box sx={{ mt: 3 }}>
                <Typography variant="h6" align="center" gutterBottom>Working Hours</Typography>
                <Typography align="center">
                  {visitingInfo?.workingHours ?? ''}
                </Typography>
              </Box>
            )}

            {propertyCondition && (
              <Box sx={{ mt: 3 }}>
                <Typography variant="h6" align="center" gutterBottom>Property Status</Typography>
                <Typography align="center">{propertyCondition}</Typography>
              </Box>
            )}
          </CardContent>
        </Card>

        <Box sx={{ mt: 3 }}>
          <Typography variant="h5" gutterBottom>Sources Of Information</Typography>
          <Typography gutterBottom>You may want to visit the following links to get more detailed data.</Typography>
          {sources && sources.length > 0 ? (
            <List>
              {sources.map((source: string, index: number) => (
                <ListItem key={index} disableGutters>
                  <ListItemIcon>
                    <InfoIcon />
                  </ListItemIcon>
                  <ListItemText>
                    <Link target="_blank" rel="noopener noreferrer" href={source}>
                      {source}
                    </Link>
                  </ListItemText>
                </ListItem>
              ))}
            </List>
          ) : (
            <Typography>No Sources Available</Typography>
          )}
        </Box>
      </Box>
    </Container>
  );
}
