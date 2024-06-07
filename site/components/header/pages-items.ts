export type MenuPageItem = {
  name: string;
  link: string;
};

export const menuItems = [
  {
    name: "Index",
    link: '/index',
  },
  {
    name: "Data Sources",
    link: '/data-sources',
  },
  {
    name: "Tech Stuff",
    link: '/tech-stuff',
  },
  {
    name: "About Us",
    link: '/about',
  },
] as MenuPageItem[];
