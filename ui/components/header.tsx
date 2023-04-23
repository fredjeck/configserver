import { PropsWithChildren } from "react";
import Link from "next/link"

interface HeaderProps {
  selected: string;
}

export default function Header({ selected }: PropsWithChildren<HeaderProps>) {
  function itemClassName(item: string) {
    const base = "mt-2 flex items-center text-sm text-gray-500";
    return item === selected ? " font-extrabold" : "";
  }

  return (
    <div className="lg:flex lg:items-center lg:justify-between p-8">
      <div className="min-w-0 flex-1">
        <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
          ConfigServer
        </h2>
        <div className="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">
          <div className={itemClassName("register")}>
            <Link href="register">Register</Link>
          </div>
          <div className={itemClassName("encrypt")}>
            <Link href="encrypt">Encrypt</Link>
          </div>
          <div className={itemClassName("statistics")}>
            <Link href="/">Statistics</Link>
          </div>
        </div>
      </div>
    </div>
  );
}
