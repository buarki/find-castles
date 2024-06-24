import { CountryCode } from "../country";
import { getDBClient } from "./client";
import { Castle } from "./model";

type GetCastleProps = {
  name: string;
  country: CountryCode;
};

export async function getCastle({ name, country }: GetCastleProps): Promise<Castle> {
  const dbClient = await getDBClient();
  const collections = dbClient.db('find-castles').collection('castles');
  return collections.findOne({
    country,
    name,
  }, {
    // matchingTags: false,
  }) as any as Castle;
}
