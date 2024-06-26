import { MetadataProps } from "@find-castles/lib/metadata-props";
import { Box, Container, Typography } from "@mui/material";
import { Metadata, ResolvingMetadata } from "next";

const siteHost = process.env.SITE_HOST;

const description = `Understand more about Find Castles`;

export async function generateMetadata(
  { params, searchParams }: MetadataProps,
  parent: ResolvingMetadata
): Promise<Metadata> {
  return {
    title: "About - Find Castles",
    description: description,
    keywords: ["castles", "heritage", "european castles", "data sources", "historical castles", "tracked countries", "untracked countries"],
    applicationName: 'Find Castles',
    robots: { index: true, follow: true },
    authors: {
      name: 'Aurelio Buarque',
      url: 'https://buarki.com'
    },
    openGraph: {
      title: "About - Find Castles",
      description: description,
      url: `${siteHost}/about`,
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
      title: "About - Find Castles",
      description: description,
      images: `${siteHost}/og.png`,
    }
  };
};

export default function AboutPage() {
  return (
    <Box sx={{ bgcolor: 'secondary.main', minHeight: '100vh', py: 4 }}>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: JSON.stringify({
          "@context": "https://schema.org",
          "@type": "ProfilePage",
          "dateCreated": "2023-06-26T00:34:00-00:00",
          "dateModified": "2024-06-26T2:59:00-00:00",
          "image": `${siteHost}/og.png`,
          "mainEntity": {
            "@type": "Person",
            "name": "Find Castles",
            "alternateName": "Find Castles",
            "description": description,
            "image": `${siteHost}/og.png`,
            "sameAs": [
              "https://twitter.com/buarki",
              "https://www.linkedin.com/in/aurelio-buarque/",
              "https://github.com/buarki",
            ]
          }
        }) }}
      />
      <Container maxWidth="md">
        <Typography variant="h1" gutterBottom>
          Why This Project Was Built? The Castles!
        </Typography>

        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, }}>
        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
            The Significance of Castles in European History
          </Typography>
          <Typography paragraph>
            With the fall of the Western Roman Empire, the Early Middle Ages began, and it culminated with the Fall of Constantinople. Undoubtedly, one of the most iconic features of this period is the proliferation of castles.
          </Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
            Capturing the Essence of European Castles
          </Typography>
          <Typography paragraph>
            This project aims to capture the essence of European castles, viewing them not merely as historical artifacts but as living repositories of culture, heritage, and human ingenuity.
          </Typography>
          <Typography>
            Through meticulous research and data aggregation, we&apos;ve embarked on a journey to consolidate information about these castles into a single platform, accessible to enthusiasts, scholars, and curious minds alike.
          </Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
          Unraveling History Through Castles
          </Typography>
          <Typography paragraph>
          Our endeavor goes beyond cataloging stone walls and towers; it&apos;s about unraveling the rich tapestry of history woven within each fortress&apos;s walls.
          </Typography>
          <Typography>From the towering bastions of medieval Portugal to the rugged keeps of Scotland, every castle holds within it tales of battles won and lost, of kings and queens, of intrigue and romance.</Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
          Democratizing Access to Castle Data
          </Typography>
          <Typography paragraph>
          Our primary objective is to democratize access to castle data across Europe, making it easily accessible for both humans and machines.
          </Typography>
          <Typography>By consolidating information about European castles into a single platform and providing an intuitive interface, we aim to break down barriers to access and empower individuals, researchers, and enthusiasts to explore and learn about these historical landmarks.</Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
          Introduction of Transparency Index
          </Typography>
          <Typography paragraph>
          We propose the introduction of a “Transparency Index“ for castle data, providing visibility into the availability and quality of information for each country in Europe.
          </Typography>
          <Typography>This index will serve as a benchmark, evaluating the comprehensiveness and reliability of castle data sources. By quantifying transparency, we can identify gaps and disparities in data accessibility, paving the way for targeted improvements.</Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
          Visibility and Accountability
          </Typography>
          <Typography paragraph>
          Through the transparency index, we seek to bring greater visibility to the state of castle data across Europe.
          </Typography>
          <Typography>By publicly highlighting areas of strength and weakness in data transparency, we aim to hold stakeholders, including governments, heritage organizations, and data custodians, accountable for the quality and accessibility of castle information within their respective countries.</Typography>
        </section>

        <section>
          <Typography variant="h2" fontSize={50} gutterBottom>
          Initiatives for Improvement
          </Typography>
          <Typography paragraph>
          Armed with insights from the transparency index, we aspire to catalyze initiatives aimed at enhancing the sources of information about European castles.
          </Typography>
          <Typography>By fostering collaboration among stakeholders, advocating for open data policies, and incentivizing data-sharing practices, we envision a future where the availability and accuracy of castle data are continually improved.</Typography>
        </section>
        </Box>
      </Container>
    </Box>
  );
}
