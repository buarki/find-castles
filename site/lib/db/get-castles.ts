import { Country, CountryCode } from "../country";
import { getDBClient } from "./client";
import { Castle } from "./model";

type GetCastlesProps = {
  contryCodes: CountryCode[];
};

export async function getCastles({ contryCodes }: GetCastlesProps): Promise<Castle[]> {
  const dbClient = await getDBClient();
  const collections = dbClient.db('find-castles').collection('castles');
  return collections.find({
    country: {
      $in: contryCodes,
    },
  }).toArray() as any as Castle[];
}
