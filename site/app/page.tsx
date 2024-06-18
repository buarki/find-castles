import { CountrySelector } from "@find-castles/components/country-selector/country-selector.component";
import { Country } from "@find-castles/lib/country";
import { Box, Button, Container, FormControl, InputLabel, MenuItem, Select, Typography } from "@mui/material";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "find castles",
  description: "find castles",
};

const countries = [
  {
    code: 'pt',
    name: 'portugal',
  },
  {
    code: 'uk',
    name: 'united kingdom',
  },
  {
    code: 'ir',
    name: 'ireland',
  },
] as Country[];

export default function Home() {
  return (
    <Box>
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
          backgroundImage: `url('/background.png')`,
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
            
            <CountrySelector countries={countries}/>
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
          backgroundImage: `url('/background-mobile.png')`,
          backgroundSize: 'cover',
          backgroundPosition: 'center',
          bgcolor: 'secondary.main',
      }}>
        <Container>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, }}>
            <Typography variant="h1" color='primary.main'>Choose A Country</Typography>
            <Typography>Embark on a journey with Castle Explorer and uncover the secrets of these timeless treasures. Select the country to see the castles and start your adventure today!</Typography>
            
            <CountrySelector countries={countries}/>
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
