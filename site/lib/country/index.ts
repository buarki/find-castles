export type CountryCode = 'at' | 'be' | 'bg' | 'hr' | 'cy' | 'cz' | 'dk' | 'ee' | 'fi' | 'fr' | 'de' | 'gr' | 'hu' | 'ie' | 'it' | 'lv' | 'lt' | 'lu' | 'mt' |'nl' | 'pl' | 'pt' | 'ro' | 'sk' | 'si' | 'es' | 'se' | 'gb'; 

export type TrackingStatus = 'tracked' | 'not-tracked';

export type Country = {
  name: string;
  code: CountryCode;
  trackingStatus: TrackingStatus;
  sources?: string[];
};

export const countries = [
  { name: "Austria", trackingStatus: 'not-tracked', code: "at" },
  { name: "Belgium", trackingStatus: 'not-tracked', code: "be" },
  { name: "Bulgaria", trackingStatus: 'not-tracked', code: "bg" },
  { name: "Croatia", trackingStatus: 'not-tracked', code: "hr" },
  { name: "Cyprus", trackingStatus: 'not-tracked', code: "cy" },
  { name: "Czech Republic", trackingStatus: 'not-tracked', code: "cz" },
  { name: "Denmark", trackingStatus: 'not-tracked', code: "dk" },
  { name: "Estonia", trackingStatus: 'not-tracked', code: "ee" },
  { name: "Finland", trackingStatus: 'not-tracked', code: "fi" },
  { name: "France", trackingStatus: 'not-tracked', code: "fr" },
  { name: "Germany", trackingStatus: 'not-tracked', code: "de" },
  { name: "Greece", trackingStatus: 'not-tracked', code: "gr" },
  { name: "Hungary", trackingStatus: 'not-tracked', code: "hu" },
  { name: "Ireland", trackingStatus: 'tracked', code: "ie" },
  { name: "Italy", trackingStatus: 'not-tracked', code: "it" },
  { name: "Latvia", trackingStatus: 'not-tracked', code: "lv" },
  { name: "Lithuania", trackingStatus: 'not-tracked', code: "lt" },
  { name: "Luxembourg", trackingStatus: 'not-tracked', code: "lu" },
  { name: "Malta", trackingStatus: 'not-tracked', code: "mt" },
  { name: "Netherlands", trackingStatus: 'not-tracked', code: "nl" },
  { name: "Poland", trackingStatus: 'not-tracked', code: "pl" },
  { name: "Portugal", trackingStatus: 'tracked', code: "pt" },
  { name: "Romania", trackingStatus: 'not-tracked', code: "ro" },
  { name: "Slovakia", trackingStatus: 'tracked', code: "sk" },
  { name: "Slovenia", trackingStatus: 'not-tracked', code: "si" },
  { name: "Spain", trackingStatus: 'not-tracked', code: "es" },
  { name: "Sweden", trackingStatus: 'not-tracked', code: "se" },
  { name: "United Kingdom", trackingStatus: 'tracked', code: "gb" },
] as Country[];
