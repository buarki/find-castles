import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Data Sources",
  description: "Find Castles Project",
};

type DataSourcePageProps = {
  params: {
    source: string;
  }
};


// what is the source?
// is the source official?
export default function DataSourcePage({ params: { source } }: DataSourcePageProps) {
  return (
    <main className="flex flex-col py-3 gap-9">
      <h1>Sources of {source}</h1>
    </main>
  );
}
