"use client";

import { Country } from "@find-castles/lib/country";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { FormControl, InputLabel, Link, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import React from "react";

export type CountrySelectorProps = {
  countries: Country[];
};

export function CountrySelector({ countries }: CountrySelectorProps) {
  const [country, setCountry] = React.useState<Country>(countries[0]);

  const handleChange = (event: SelectChangeEvent<string>, _child: React.ReactNode) => {
    setCountry(countries.find((c) => c.code === event.target.value)!);
  };

  return (
   <FormControl
      fullWidth
      sx={{
        display: 'flex',
        flexDirection: 'column',
        gap: 3,
      }}>
      <InputLabel id="demo-simple-select-label">Country</InputLabel>
      <Select
        labelId="demo-simple-select-label"
        id="demo-simple-select"
        value={country.code}
        label="Country"
        onChange={handleChange}
      >
        {
          countries.map((country: Country) => (
            <MenuItem key={country.code} value={country.code}>{ toTitleCase(country.name) }</MenuItem>
          ))
        }
      </Select>
      <Link
        href={`/castles/?country=${country.code}`}
        textAlign='center'
        sx={{
          bgcolor: 'primary.main',
          color: 'secondary.main',
          width: '100%',
          py: 1,
        }}>
          Search
      </Link>
    </FormControl>
  );
} 
