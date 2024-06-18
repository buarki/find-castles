import { Inter } from "next/font/google";
import "./globals.css";
import { Box } from "@mui/material";
import { ThemeRegistry } from "@find-castles/theme/theme-registry";
import { Header } from "@find-castles/components/header/header.component";
import { Footer } from "@find-castles/components/footer/footer.component";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <ThemeRegistry options={{ key: 'mui-theme' }}>
          <Box sx={{
            display: 'flex',
            flexDirection: 'column',
            minHeight: '100vh',
          }}>
            <Header />
            <Box sx={{ flexGrow: 1}}>
              {children}
            </Box>
            <Footer />
          </Box>
        </ThemeRegistry>
      </body>
    </html>
  );
}