import { MidAndLargeScreeMenu } from "./mid-and-large-scree-menu.component";
import { menuItems } from "./pages-items";
import { SmallScreenMenu } from "./small-screen-menu.component";

type HeaderProps = {
  className?: string;
};

export function Header({ className }: HeaderProps) {
  return (
    <header
      className={`
        ${className}
        w-full
        flex
        justify-between
        p-6
        border-b-2
      `}>
      <a href="/">Find Castles</a>
      <MidAndLargeScreeMenu menuItems={menuItems} className="md:flex hidden"/>
      <SmallScreenMenu menuItems={menuItems} className="md:hidden"/>
    </header>
  );
}
