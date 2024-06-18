"use client";

import { Box, IconButton, Link, Menu, MenuItem, Typography } from "@mui/material";
import MenuIcon from '@mui/icons-material/Menu';
import React from "react";
import { MenuItemData } from "./menu-item";
import { toTitleCase } from "@find-castles/lib/to-title-case";

type SmallScreenMenuProps = {
  menuItems: MenuItemData[];
  sx?: any;
};

export function SmallScreenMenu({ menuItems, sx }: SmallScreenMenuProps) {
  const [anchorElNav, setAnchorElNav] = React.useState<null | HTMLElement>(null);

  const handleOpenNavMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElNav(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(null);
  };

  return (
    <Box
      sx={{
        ...sx,
      }}>
      <IconButton
        onClick={handleOpenNavMenu}
        edge="end"
        color="inherit"
        aria-label="menu"
        >
        <MenuIcon sx={{ fontSize: '2.5rem' }}/>
      </IconButton>
      <Menu
        anchorEl={anchorElNav}
        onClose={handleCloseNavMenu}
        id="menu-appbar"
        open={Boolean(anchorElNav)}>
        {
          menuItems.map((menuItem: MenuItemData, index: number) => (
            <MenuItem
              sx={{
                m: 0,
                p: 0,
              }}
              key={index}>
              <Link
                sx={{
                  display: 'block',
                  width: '100%',
                  height: '100%',
                  padding: '8px 16px',
                }}
                underline="none"
                href={menuItem.href}>
                <Typography
                  textAlign='center'
                  variant="h3"
                  >
                    { toTitleCase(menuItem.name) }
                  </Typography>
              </Link>
            </MenuItem>
          ))
        }
      </Menu>
    </Box>
  );
}
