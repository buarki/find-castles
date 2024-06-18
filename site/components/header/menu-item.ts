export type MenuItemData = {
  name: string;
  href: string;
};

export const menuItems = [
  {
    name: 'about',
    href: '/about',
  },
  {
    name: 'castles',
    href: '/castles',
  },
  {
    name: 'data sources',
    href: '/data-sources',
  },
  {
    name: 'tech stuff',
    href: '/tech-stuff',
  },
] as MenuItemData[];

