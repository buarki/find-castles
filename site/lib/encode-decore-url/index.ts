import { CountryCode } from "../country";

const spaceReplacer = "-";

export function encodeCastleURL(castleName: string, countryCode: CountryCode): string {
  return encodeURIComponent(`${castleName.replaceAll(" ", spaceReplacer)}${spaceReplacer}${countryCode}`);
}

export function encodeCountryCastleListURL(countryName: string): string {
  return encodeURIComponent(countryName.replaceAll(" ", spaceReplacer));
}


export function decodeCastleURL(rawURL: string): { castleName: string, countryCode: CountryCode } {
  const tokens = decodeURIComponent(rawURL).split(spaceReplacer);
  const countryCode = tokens[tokens.length - 1] as CountryCode;
  const castleName = tokens.slice(0, tokens.length-1).join(" ");
  return {
    countryCode,
    castleName,
  };
}
