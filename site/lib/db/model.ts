import { CountryCode } from "../country";

export interface Contact {
  phone?: string;
  email?: string;
}

export interface Facilities {
  assistanceDogsAllowed: boolean;
  cafe: boolean;
  restrooms: boolean;
  giftshops: boolean;
  pinicArea: boolean;
  parking: boolean;
  exhibitions: boolean;
  wheelchairSupport: boolean;
}

export interface VisitingInfo {
  workingHours: string;
  facilities?: Facilities;
}

export interface Castle {
  _id: string;
  country: CountryCode;
  name: string;
  city: string;
  contact?: Contact;
  coordinates: string;
  pictureURL: string;
  sources: string[];
  state: string;
  district?: string;
  visitingInfo?: VisitingInfo;
  webName: string;
}
