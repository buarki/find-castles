"use client";

import { useState } from "react";
import { AiOutlineMenu } from "react-icons/ai";
import { MenuPageItem } from "./pages-items";

type SmallScreenMenuProps = {
  menuItems: MenuPageItem[];
  className?: string;
};

export function SmallScreenMenu({ menuItems, className }: SmallScreenMenuProps) {
  const [showMenu, setShowMenu] = useState(false);
  
  const showMenuEvent = () => {
    setShowMenu(!showMenu);
  }
  
  return (
    <div
      className={`
        relative
        w-full
        flex
        justify-end
        ${className}
      `}>
      <button aria-label="Menu" onClick={() => showMenuEvent()}>
        <AiOutlineMenu className="text-3xl"/>
      </button>
      {
        showMenu &&
        <div className="absolute z-10 divide-y rounded-lg shadow w-44 bg-gray-900 mt-10">
          <ul className="py-2 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownDefaultButton">
            <li>
              <a href="/" className="block px-4 py-2 text-white">Home</a>
            </li>
            <li>
              <a href="/about" className="block px-4 py-2 text-white">About</a>
            </li>
            <li>
              <a href="/blog" className="block px-4 py-2 text-white">Blog</a>
            </li>
            <li>
              <a href="/portfolio" className="block px-4 py-2 text-white">Portfolio</a>
            </li>
            <li>
              <a href="/appearances" className="block px-4 py-2 text-white">Appearances</a>
            </li>
          </ul>
        </div>
      }
    </div>
  );
}