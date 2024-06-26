"use client";

import { Country, CountryCode } from "@find-castles/lib/country";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Grid, CircularProgress } from "@mui/material";
import React, { useState } from "react";
import axios from "axios";
import { CastleCard } from "./castle-card.component";
import { Castle } from "@find-castles/lib/db/model";
import Link from "next/link";

export type ClientSideCastlesListerProps = {
  countries: Country[];
  currentCountry?: CountryCode;
};

export function ClientSideCastlesLister({ countries, currentCountry }: ClientSideCastlesListerProps) {
  countries = [{ name: "Select", code: 'fake' as CountryCode } as Country, ...countries];
  const [country, setCountry] = useState<Country>(currentCountry ? countries.find((c) => c.code === currentCountry)! : countries[0]);
  const [castles, setCastles] = useState<Castle[]>([]);
  const [loading, setLoading] = useState(false);
  const [availableStates, setAvailableStates] = useState<string[]>([]);
  const [availablePropertyConditions, setAvailablePropertyConditions] = useState<string[]>([]);
  const [selectedState, setSelectedState] = useState<string | undefined>(undefined);
  const [selectedPropertyCondition, setSelectedPropertyCondition] = useState<string | undefined>(undefined);

  const handleChange = async (event: SelectChangeEvent<string>, _child: React.ReactNode) => {
    const selectedCountry = countries.find((c) => c.code === event.target.value)!;
    setCountry(selectedCountry);

    setLoading(true);
    try {
      const response = await axios.get(`/castles/api/?country=${selectedCountry.code}`);
      const fetchedCastles = response.data.data as Castle[];
      setSelectedState(undefined);
      setSelectedPropertyCondition(undefined);

      setCastles(fetchedCastles);

      const states = Array.from(new Set(fetchedCastles.map(castle => castle.state)));
      const propertyConditions = Array.from(new Set(fetchedCastles.map(castle => castle.propertyCondition)));

      setAvailableStates(states);
      setAvailablePropertyConditions(propertyConditions);
    } catch (error) {
      console.error("Error fetching castles:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleStateChange = (event: SelectChangeEvent<string>) => {
    setSelectedState(event.target.value);
  };

  const handlePropertyConditionChange = (event: SelectChangeEvent<string>) => {
    setSelectedPropertyCondition(event.target.value);
  };

  const applyFilters = (castle: Castle) => {
    let showCastle = true;

    if (selectedState && castle.state !== selectedState) {
      showCastle = false;
    }

    if (selectedPropertyCondition && castle.propertyCondition !== selectedPropertyCondition) {
      showCastle = false;
    }

    return showCastle;
  };

  const filteredCastles = castles.filter(applyFilters);

  return (
    <Box sx={{ my: 6 }}>
      <FormControl fullWidth sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
        <InputLabel id="demo-simple-select-label">Country</InputLabel>
        <Select
          labelId="demo-simple-select-label"
          id="demo-simple-select"
          value={country.code}
          label="Country"
          onChange={handleChange}
        >
          {countries.map((country: Country) => (
            <MenuItem key={country.code} value={country.code}>
              {toTitleCase(country.name)}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          {availableStates.length > 0 && (
            <FormControl fullWidth sx={{ display: 'flex', flexDirection: 'column', gap: 3, mt: 3 }}>
              <InputLabel id="state-label">State</InputLabel>
              <Select
                label="State"
                labelId="state-label"
                id="state-select"
                value={selectedState || ""}
                onChange={handleStateChange}
              >
                <MenuItem value="">All</MenuItem>
                {availableStates.map(state => (
                  <MenuItem key={state} value={state}>{state}</MenuItem>
                ))}
              </Select>
            </FormControl>
          )}

          {availablePropertyConditions.length > 0 && (
            <FormControl fullWidth sx={{ display: 'flex', flexDirection: 'column', gap: 3, mt: 3 }}>
              <InputLabel id="property-condition-label">Property Condition</InputLabel>
              <Select
                label="Property Condition"
                labelId="property-condition-label"
                id="property-condition-select"
                value={selectedPropertyCondition || ""}
                onChange={handlePropertyConditionChange}
              >
                <MenuItem value="">All</MenuItem>
                {availablePropertyConditions.map(condition => (
                  <MenuItem key={condition} value={condition}>{condition}</MenuItem>
                ))}
              </Select>
            </FormControl>
          )}

          <Grid container spacing={3} sx={{ mt: 3 }}>
            {filteredCastles.map((castle) => (
              <Grid item key={castle._id} xs={12} sm={6} md={4}>
                <Link href={`/castles/${castle.webName}`}>
                  <CastleCard castle={castle} />
                </Link>
              </Grid>
            ))}
          </Grid>
        </>
      )}
    </Box>
  );
}
