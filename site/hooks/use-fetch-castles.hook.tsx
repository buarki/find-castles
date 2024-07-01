'use-client';

import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { Castle } from '@find-castles/lib/db/model';
import { CountryCode } from '@find-castles/lib/country';

export function useFetchCastles(countryCode: CountryCode) {
  const [castles, setCastles] = useState<Castle[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | undefined>(undefined);

  const fetchCastles = useCallback(async (countryCode: CountryCode) => {
    setLoading(true);
    try {
      const response = await axios.get(`/castles/api/?country=${countryCode}`);
      setCastles(response.data.data);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (countryCode) {
      fetchCastles(countryCode);
    }
  }, [countryCode, fetchCastles]);

  return { castles, loading, error };
}

