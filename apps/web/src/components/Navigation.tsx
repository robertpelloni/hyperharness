"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Sheet, SheetContent, SheetTrigger } from "@borg/ui";
import { Button } from "@borg/ui";
import { Menu } from "lucide-react";
import { useState } from "react";

const NAV_ITEMS = [
    { href: '/', label: 'Mission Control', color: 'hover:text-blue-500', activeColor: 'text-blue-500' },
    { href: '/docs', label: 'Documentation', color: 'hover:text-blue-500', activeColor: 'text-blue-500' },
    { href: '/status', label: 'System Status', color: 'hover:text-blue-500', activeColor: 'text-blue-500' },
    { href: '/dashboard/council', label: 'Council', color: 'hover:text-purple-500', activeColor: 'text-purple-500' },
    { href: '/dashboard/skills', label: 'Skills', color: 'hover:text-green-500', activeColor: 'text-green-500' },
    { href: '/dashboard/reader', label: 'Reader', color: 'hover:text-orange-500', activeColor: 'text-orange-500' },
    { href: '/dashboard/command', label: 'Command Center', color: 'hover:text-red-500', activeColor: 'text-red-500' },
    { href: '/dashboard/inspector', label: 'Traffic', color: 'hover:text-yellow-500', activeColor: 'text-yellow-500' },
    { href: '/dashboard/config', label: 'Settings', color: 'hover:text-slate-500', activeColor: 'text-slate-500' },
];

export function Navigation() {
    const pathname = usePathname();
    const [open, setOpen] = useState(false);

    const isActive = (path: string) => pathname === path;

    return (
        <nav className="w-full bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800 px-6 py-4 flex items-center justify-between sticky top-0 z-50">
            <div className="flex items-center gap-6">
                <div className="text-xl font-bold bg-gradient-to-r from-blue-500 to-purple-500 bg-clip-text text-transparent">
                    BORG
                </div>

                {/* Desktop Navigation */}
                <div className="hidden md:flex gap-4">
                    {NAV_ITEMS.map((item) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`text-sm font-medium transition-colors ${item.color} ${isActive(item.href) ? item.activeColor : 'text-zinc-500 dark:text-zinc-400'}`}
                        >
                            {item.label}
                        </Link>
                    ))}
                </div>
            </div>

            {/* Mobile Navigation */}
            <div className="md:hidden">
                <Sheet open={open} onOpenChange={setOpen}>
                    <SheetTrigger asChild>
                        <Button variant="ghost" size="icon">
                            <Menu className="h-6 w-6" />
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-[300px] sm:w-[400px]">
                        <div className="flex flex-col gap-4 mt-8">
                            {NAV_ITEMS.map((item) => (
                                <Link
                                    key={item.href}
                                    href={item.href}
                                    onClick={() => setOpen(false)}
                                    className={`text-lg font-medium transition-colors ${item.color} ${isActive(item.href) ? item.activeColor : 'text-zinc-500 dark:text-zinc-400'}`}
                                >
                                    {item.label}
                                </Link>
                            ))}
                            <div className="mt-auto pt-8 border-t border-zinc-200 dark:border-zinc-800">
                                <span className="text-xs text-zinc-400">v0.1.0-alpha</span>
                            </div>
                        </div>
                    </SheetContent>
                </Sheet>
            </div>

            <div className="hidden md:block text-xs text-zinc-400">
                v0.1.0-alpha
            </div>
        </nav>
    );
}
