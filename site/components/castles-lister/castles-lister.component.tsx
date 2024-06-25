"use client";

import { Country, CountryCode } from "@find-castles/lib/country";
import { toTitleCase } from "@find-castles/lib/to-title-case";
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent, Grid, CircularProgress } from "@mui/material";
import React, { useEffect, useState } from "react";
import axios from "axios";
import { CastleCard } from "./castle-card.component";
import { Castle } from "@find-castles/lib/db/model";
import Link from "next/link";

export type ClientSideCastlesListerProps = {
  countries: Country[];
  currentCountry?: CountryCode;
};

export function ClientSideCastlesLister({ countries, currentCountry }: ClientSideCastlesListerProps) {
  console.log('mandou', currentCountry);
  const [country, setCountry] = useState<Country>(currentCountry ? countries.find((c) => c.code === currentCountry)! : countries[0]);
  const [castles, setCastles] = useState<Castle[]>([]);
  const [loading, setLoading] = useState(false);

  const handleChange = (event: SelectChangeEvent<string>, _child: React.ReactNode) => {
    const selectedCountry = countries.find((c) => c.code === event.target.value)!;
    setCountry(selectedCountry);
  };

  useEffect(() => {
    const fetchCastles = async () => {
      setLoading(true);
      try {
        const response = await axios.get(`/castles/api/?country=${country.code}`);
        setCastles(response.data.data);
      } catch (error) {
        console.error("Error fetching castles:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchCastles();
  }, [country]);

  return (
    <Box
      sx={{
        my: 6,
      }}
      >
      <FormControl
        fullWidth
        sx={{
          display: 'flex',
          flexDirection: 'column',
          gap: 3,
        }}
      >
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
        <Grid
            container
            spacing={3}
            sx={{ mt: 3 }
            }
          >
          {castles.map((castle) => (
            <Grid item key={castle._id} xs={12} sm={6} md={4}>
              <Link href={`/castles/${castle.webName}`}>
                <CastleCard castle={castle} />
              </Link>
            </Grid>
          ))}
        </Grid>
      )}
    </Box>
  );
}
