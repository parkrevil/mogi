'use client'

import { navigation } from '@/libs/navigation';
import { usePathname } from "next/navigation";

export default function Header() {
  const pathname = usePathname();
  console.log(pathname);
  const activeMenu = navigation.find((item) => item.href === pathname);

  return (
    <header className="mt-6 mx-7">
      <h2 className="text-2xl font-bold">{activeMenu?.label}</h2>
    </header>
  )
}
