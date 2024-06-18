import { AppBar, Box, Container, Link, MenuItem, Toolbar, Typography } from "@mui/material";
import Image from "next/image";
import { SmallScreenMenu } from "./small-screen.component";
import { menuItems, MenuItemData} from "./menu-item";
import { toTitleCase } from "@find-castles/lib/to-title-case";

export function Header() {
  return (
    <AppBar
      position="sticky"
      color="secondary"
      sx={{ boxShadow: 'none', zIndex: 1, }}>
      <Container maxWidth="lg">
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
            alignItems: 'center'
          }}>
          <Link href="/">
            <Image alt="Logo" width={100} height={40} src={'/logo.png'}/>
          </Link>
          <Toolbar>
            <Box sx={{
              display: {
                xs: 'none',
                md: 'flex',
              },
              gap: 2,
              alignItems: 'center'
            }}>
              {
                menuItems.map((menuItem: MenuItemData, index: number) => (
                  <Link
                    underline="none"
                    href={menuItem.href}
                    key={index}
                    sx={{
                      p: 1,
                      ':hover': {
                        bgcolor: 'secondary.dark'
                      }
                    }}
                    >
                    <Typography variant="h3">{ toTitleCase(menuItem.name) }</Typography>
                  </Link>
                ))
              }
            </Box>
            <SmallScreenMenu
              sx={{
                display: {
                  xs: 'visible',
                  md: 'none',
                },
              }}
              menuItems={menuItems}/>
          </Toolbar>
        </Box>
      </Container>
    </AppBar>
  );
}
