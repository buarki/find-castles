import React from 'react';
import { Box, Card, CardContent, CardMedia, Typography } from '@mui/material';
import { Castle } from '@find-castles/lib/db/model';
import { toTitleCase } from '@find-castles/lib/to-title-case';

type CastleCardProps = {
  castle: Castle;
};

export function CastleCard({ castle }:CastleCardProps) {
  const { pictureURL, name, city, visitingInfo } = castle;
  return (
    <Card
      sx={{
        width: { xs: '100%', sm: '100%', md: '100%' },
        aspectRatio: '1',
        transition: 'transform 0.3s ease-in-out',
        '&:hover': {
          transform: 'scale(1.05)',
        },
      }}
      >
      <Box sx={{ position: 'relative', height: '100%', width: '100%' }}>
        <CardMedia
          component="img"
          image={pictureURL.includes("https://") ? pictureURL : `https://${pictureURL}`}
          alt={name}
          sx={{ position: 'absolute', top: 0, left: 0, height: '100%', width: '100%', objectFit: 'cover' }}
        />
        <CardContent
          sx={{
            position: 'absolute',
            bottom: 0,
            left: 0,
            width: '100%',
            background: 'rgba(0, 0, 0, 0.5)',
            color: 'white',
          }}
        >
          <Typography variant="h6">{toTitleCase(name)}</Typography>
          <Typography variant="body2">{toTitleCase(city)}</Typography>
        </CardContent>
      </Box>
    </Card>
  );
};

