import "./globals.css";
import { Roboto } from '@next/font/google';
import { Header } from "@find-castles/components/header/header.component";
import { Footer } from "@find-castles/components/footer/footer.component";

const roboto = Roboto({
  weight: ['400', '700'],
  subsets: ['latin'],
  variable: '--font-roboto',
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="flex flex-col min-h-screen">
        <Header className="px-12"/>
        <div className="grow px-12">
          {children}
        </div>
        <Footer className="px-12"/>
      </body>
    </html>
  );
}
