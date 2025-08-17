'use client'

import { navigation } from '@/libs/navigation';
import { usePathname } from "next/navigation";

export default function Header() {
  const pathname = usePathname();
  const activeMenu = navigation.find((item) => item.href === pathname);

  return (
    <header>
      <div className="flex flex-row items-center text-2xl font-bold">
        <div className="flex justify-center items-center w-8 mr-3 text-center -ml-1">{activeMenu?.icon}</div>
        <h2>{activeMenu?.label}</h2>
      </div>
    </header>
  )
}
