import { getDBClient } from "./client";
import { Castle } from "./model";


export async function getCastle(webName: string): Promise<Castle> {
  const dbClient = await getDBClient();
  const collections = dbClient.db('find-castles').collection('castles');
  return collections.findOne({
    webName,
  }) as any as Castle;
}
