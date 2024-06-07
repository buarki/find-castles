import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Find Castles - About",
  description: "Find Castles Project",
};

export default function AboutPage() {
  return (
    <main className="flex flex-col justify-center items-center">
      <article className="container mx-auto py-8 px-4 w-8/12">
        <h1 className="text-3xl font-bold mb-4">Why This Project Was Built? The Castles!</h1>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">The Significance of Castles in European History</h2>
          <p>With the fall of the Western Roman Empire, the Early Middle Ages began, and it culminated with the Fall of Constantinople. Undoubtedly, one of the most iconic features of this period is the proliferation of castles.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Capturing the Essence of European Castles</h2>
          <p>This project aims to capture the essence of European castles, viewing them not merely as historical artifacts but as living repositories of culture, heritage, and human ingenuity.</p>
          <p>Through meticulous research and data aggregation, we&apos;ve embarked on a journey to consolidate information about these castles into a single platform, accessible to enthusiasts, scholars, and curious minds alike.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Unraveling History Through Castles</h2>
          <p>Our endeavor goes beyond cataloging stone walls and towers; it&apos;s about unraveling the rich tapestry of history woven within each fortress&apos;s walls.</p>
          <p>From the towering bastions of medieval Portugal to the rugged keeps of Scotland, every castle holds within it tales of battles won and lost, of kings and queens, of intrigue and romance.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Democratizing Access to Castle Data</h2>
          <p>Our primary objective is to democratize access to castle data across Europe, making it easily accessible for both humans and machines.</p>
          <p>By consolidating information about European castles into a single platform and providing an intuitive interface, we aim to break down barriers to access and empower individuals, researchers, and enthusiasts to explore and learn about these historical landmarks.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Introduction of Transparency Index</h2>
          <p>We propose the introduction of a &ldquo;Transparency Index&ldquo; for castle data, providing visibility into the availability and quality of information for each country in Europe.</p>
          <p>This index will serve as a benchmark, evaluating the comprehensiveness and reliability of castle data sources. By quantifying transparency, we can identify gaps and disparities in data accessibility, paving the way for targeted improvements.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Visibility and Accountability</h2>
          <p>Through the transparency index, we seek to bring greater visibility to the state of castle data across Europe.</p>
          <p>By publicly highlighting areas of strength and weakness in data transparency, we aim to hold stakeholders, including governments, heritage organizations, and data custodians, accountable for the quality and accessibility of castle information within their respective countries.</p>
        </section>

        <section className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Initiatives for Improvement</h2>
          <p>Armed with insights from the transparency index, we aspire to catalyze initiatives aimed at enhancing the sources of information about European castles.</p>
          <p>By fostering collaboration among stakeholders, advocating for open data policies, and incentivizing data-sharing practices, we envision a future where the availability and accuracy of castle data are continually improved.</p>
        </section>
      </article>
    </main>
  );
}
