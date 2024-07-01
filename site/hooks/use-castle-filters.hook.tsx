'use-client';

import { useState, useEffect } from 'react';
import { Castle } from '@find-castles/lib/db/model';

export function useCastleFilters(castles: Castle[]) {
  const [availableStates, setAvailableStates] = useState<string[]>([]);
  const [availablePropertyConditions, setAvailablePropertyConditions] = useState<string[]>([]);
  const [selectedState, setSelectedState] = useState<string | undefined>(undefined);
  const [selectedPropertyCondition, setSelectedPropertyCondition] = useState<string | undefined>(undefined);

  useEffect(() => {
    const states = Array.from(new Set(castles.map(castle => castle.state)));
    const propertyConditions = Array.from(new Set(castles.map(castle => castle.propertyCondition)));

    setAvailableStates(states);
    setAvailablePropertyConditions(propertyConditions);
  }, [castles]);

  const filteredCastles = castles.filter((castle) => {
    return (!selectedState || castle.state === selectedState) &&
           (!selectedPropertyCondition || castle.propertyCondition === selectedPropertyCondition);
  });

  return {
    availableStates,
    availablePropertyConditions,
    selectedState,
    setSelectedState,
    selectedPropertyCondition,
    setSelectedPropertyCondition,
    filteredCastles
  };
}
