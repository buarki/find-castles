import { MongoClient } from 'mongodb';

const dbURL = process.env.DB_URL;

const options = {};

let dbClient: MongoClient;

// TODO check if it is really recreating it
export async function getDBClient() {
  if (!dbURL) {
    throw new Error('define DB_URL env var');
  }
  if (dbClient == null) {
    console.log('>>> CREATING DB FROM SCRATCH');
    const client = new MongoClient(dbURL, options);
    dbClient = await client.connect();
  }
  return dbClient;
}
