import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Data Sources",
  description: "Find Castles Project",
};

enum Status {
  unknown = "unknown",
  tracked = "tracked",
}

type Country = {
  name: string;
  status: Status;
  code: string;
};

const countries = [
  { name: "Austria", status: Status.unknown, code: "AT" },
  { name: "Belgium", status: Status.unknown, code: "BE" },
  { name: "Bulgaria", status: Status.unknown, code: "BG" },
  { name: "Croatia", status: Status.unknown, code: "HR" },
  { name: "Cyprus", status: Status.unknown, code: "CY" },
  { name: "Czech Republic", status: Status.unknown, code: "CZ" },
  { name: "Denmark", status: Status.unknown, code: "DK" },
  { name: "Estonia", status: Status.unknown, code: "EE" },
  { name: "Finland", status: Status.unknown, code: "FI" },
  { name: "France", status: Status.unknown, code: "FR" },
  { name: "Germany", status: Status.unknown, code: "DE" },
  { name: "Greece", status: Status.unknown, code: "GR" },
  { name: "Hungary", status: Status.unknown, code: "HU" },
  { name: "Ireland", status: Status.tracked, code: "IE" },
  { name: "Italy", status: Status.unknown, code: "IT" },
  { name: "Latvia", status: Status.unknown, code: "LV" },
  { name: "Lithuania", status: Status.unknown, code: "LT" },
  { name: "Luxembourg", status: Status.unknown, code: "LU" },
  { name: "Malta", status: Status.unknown, code: "MT" },
  { name: "Netherlands", status: Status.unknown, code: "NL" },
  { name: "Poland", status: Status.unknown, code: "PL" },
  { name: "Portugal", status: Status.tracked, code: "PT" },
  { name: "Romania", status: Status.unknown, code: "RO" },
  { name: "Slovakia", status: Status.unknown, code: "SK" },
  { name: "Slovenia", status: Status.unknown, code: "SI" },
  { name: "Spain", status: Status.unknown, code: "ES" },
  { name: "Sweden", status: Status.unknown, code: "SE" },
  { name: "United Kingdom", status: Status.tracked, code: "GB" }
] as Country[];

export default function DataSourcesPage() {
  return (
    <main className="flex flex-col py-3 gap-9">
      <h1 className="text-3xl">Data Sources</h1>
      <section className="flex flex-col gap-3">
        <header>
          <h2 className="text-xl">Countries Tracked by Us</h2>
        </header>
        <article>
          <ul>
            {
              countries.filter((country: Country) => country.status === Status.tracked).map((country, index) => (
                <li key={index}><a className="underline" href={`/data-sources/${country.code.toLowerCase()}`}>{country.name}</a></li>
              ))
            }
          </ul>
        </article>
      </section>
      <section className="flex flex-col gap-3">
        <header>
          <h2 className="text-xl">Countries NOT Tracked by Us</h2>
        </header>
        <article>
          <ul>
            {countries.filter((country: Country) => country.status !== Status.tracked).map((country, index) => (
              <li key={index}>{country.name}</li>
            ))}
          </ul>
        </article>
      </section>
    </main>
  );
}
