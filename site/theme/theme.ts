import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: "#5C3200",
    },
    secondary: {
      main: '#F9DEB8',
      dark: '#E6C896',
    },
  },
  typography: {
    fontFamily: [
      'Manrope',
      'Anton',
      'Arial',
      'sans-serif',
    ].join(','),
    h1: {
      fontSize: '3rem',
      fontWeight: 600,
      fontFamily: 'Anton, sans-serif',
    },
    h3: {
      fontSize: '1rem',
      fontWeight: 600,
      fontFamily: 'Anton, sans-serif',
    },
  },
});
