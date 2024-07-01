"use client";

import { Country, CountryCode } from "@find-castles/lib/country";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Grid, CircularProgress } from "@mui/material";
import React, { useState } from "react";
import { CastleCard } from "./castle-card.component";
import Link from "next/link";
import { useCastleFilters } from "@find-castles/hooks/use-castle-filters.hook";
import { useFetchCastles } from "@find-castles/hooks/use-fetch-castles.hook";

export type ClientSideCastlesListerProps = {
  countries: Country[];
  currentCountry?: CountryCode;
};

export function ClientSideCastlesLister({ countries, currentCountry }: ClientSideCastlesListerProps) {
  countries = [{ name: "Select", code: 'fake' as CountryCode } as Country, ...countries];
  const [country, setCountry] = useState<Country>(currentCountry ? countries.find((c) => c.code === currentCountry)! : countries[0]);
  const { castles, loading, error } = useFetchCastles(country.code);
  const {
    availableStates,
    availablePropertyConditions,
    selectedState,
    setSelectedState,
    selectedPropertyCondition,
    setSelectedPropertyCondition,
    filteredCastles
  } = useCastleFilters(castles);

  const handleCountryChange = (event: SelectChangeEvent<string>, _child: React.ReactNode) => {
    const selectedCountry = countries.find((c) => c.code === event.target.value)!;
    setCountry(selectedCountry);
  };

  const handleStateChange = (event: SelectChangeEvent<string>) => {
    setSelectedState(event.target.value);
  };

  const handlePropertyConditionChange = (event: SelectChangeEvent<string>) => {
    setSelectedPropertyCondition(event.target.value);
  };

  return (
    <Box sx={{ my: 6 }}>
      <FormControl fullWidth sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
        <InputLabel id="country-select-label">Country</InputLabel>
        <Select
          labelId="country-select-label"
          id="country-select"
          value={country.code}
          label="Country"
          onChange={handleCountryChange}
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
              <InputLabel id="state-select-label">State</InputLabel>
              <Select
                label="State"
                labelId="state-select-label"
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
              <InputLabel id="property-condition-select-label">Property Condition</InputLabel>
              <Select
                label="Property Condition"
                labelId="property-condition-select-label"
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
