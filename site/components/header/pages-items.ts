export type MenuPageItem = {
  name: string;
  link: string;
};

export const menuItems = [
  {
    name: "About",
    link: '/about',
  },
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
    name: "Crew",
    link: '/crew',
  },
] as MenuPageItem[];
