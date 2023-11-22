export type NavbarProps = {
    title: string
}

export const Navbar = ({title}: NavbarProps) => {
    return (
        <div className="navbar w-full">
        <div className="flex-none lg:hidden">
            <label htmlFor="my-drawer-3" aria-label="open sidebar" className="btn btn-square btn-ghost">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="inline-block h-6 w-6 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
            </label>
        </div>
        <div className="mx-2 flex-1 px-2 text-xl">{title}</div>
        <ul className="menu menu-horizontal gap-2">
            <input type="text" placeholder="Search" className="input input-bordered input-sm m-auto w-20 max-w-xs sm:w-28 md:w-56" />
            <label tabIndex={0} className="avatar btn btn-circle btn-ghost m-auto">
                <div className="w-10 rounded-full">
                    <img alt="Tailwind CSS Navbar component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                </div>
            </label>
        </ul>
    </div>
    )
}