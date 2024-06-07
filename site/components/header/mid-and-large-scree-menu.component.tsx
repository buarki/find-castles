import { MenuPageItem } from "./pages-items";

type MenuIconProps = {
  href: string;
  name: string;
};

function MenuIcon({ href, name  }: MenuIconProps) {
  return (
    <a
      className={`p-3 hover:bg-gray-300`}
      href={href}>
      { name }
    </a>
  );
}

type MidAndLargeScreeMenuProsp = {
  menuItems: MenuPageItem[];
  className?: string;
};

export function MidAndLargeScreeMenu({ menuItems, className }: MidAndLargeScreeMenuProsp) {
  return (
    <ul className={`${className} flex gap-3 flex items-center justify-center`}>
      {
        menuItems.map((item, index) => (<li key={index}><MenuIcon name={item.name} href={item.link}/></li>)) }
    </ul>
  );
}
