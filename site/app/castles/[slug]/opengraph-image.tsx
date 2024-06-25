import { countries } from '@find-castles/lib/country';
import { getCastle } from '@find-castles/lib/db/get-castle';
import { toTitleCase } from '@find-castles/lib/to-title-case';
import { notFound } from 'next/navigation';
import { ImageResponse } from 'next/og';

export const runtime = 'edge';

export const size = {
  width: 1200,
  height: 630,
};

export default async function Image({ params }: { params: { slug: string } }) {
  const castleWebName = params.slug;
  const foundCastle = await getCastle(castleWebName);
  if (!foundCastle) {
    notFound();
  }

  const { name, country, pictureURL } = foundCastle;

  return new ImageResponse(
    (
      <div
        style={{
          position: 'relative',
          width: '100%',
          height: '100%',
          overflow: 'hidden',
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <img
          src={pictureURL}
          alt={name}
          style={{
            flex: '1',
            width: '100%',
            height: '100%',
            objectFit: 'cover',
          }}
        />
        <div
          style={{
            background: 'rgba(0, 0, 0, 0.5)',
            color: 'white',
            padding: '10px',
            boxSizing: 'border-box',
            textAlign: 'center',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 16, }}>
            <p style={{ fontSize: '3.5rem', fontWeight: 'bold', margin: 0, }}>
              {toTitleCase(name)}
            </p>
            <p style={{ fontSize: '2.5rem', marginLeft: '10px', margin: 0, }}>
              {countries.find((c) => c.code === country)?.name}
            </p>
          </div>
          <p style={{ fontSize: '2.5rem', margin: 0 }}>Find Castles</p>
        </div>
      </div>
    ),
    { ...size }
  );
}