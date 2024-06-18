import { Typography } from "@mui/material";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Not Found",
  description: "find castles",
};

export default function NotFoundPage() {
  return (
    <Typography variant="h1">This Is Not The Path You Are Looking For</Typography>
  );
}
