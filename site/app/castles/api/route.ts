import { CountryCode } from "@find-castles/lib/country";
import { getCastles } from "@find-castles/lib/db/get-castles";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const country = url.searchParams.get('country');
  if (!country) {
    return Response.json({ message: 'missing country id' }, {status: 400});
  }
  console.log({ country });
  const foundCastles = await getCastles({ contryCodes: [country as CountryCode] });
  return Response.json({
    data: foundCastles,
  });
}
