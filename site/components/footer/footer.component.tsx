import { Box, Button, Container, Link, Typography } from "@mui/material";
import GitHubIcon from '@mui/icons-material/GitHub';

export function Footer() {
  return (
    <Box sx={{ bgcolor: 'secondary.main', py: 3 }}>
      <Container>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Box>
              <Link target="_blank" rel="noopener" title="Visit Project on Github" href="https://github.com/buarki/find-castles"><GitHubIcon /></Link>
            </Box>
            {/* <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Button>Your privacy choices</Button>
            </Box> */}
          </Box>
          <Box>
            <Typography>
              Made with love by <Link href="https://buarki.com" target="_blank" rel="noopener">buarki.com</Link>
            </Typography>
          </Box>
        </Box>
      </Container>
    </Box>
  );
}
