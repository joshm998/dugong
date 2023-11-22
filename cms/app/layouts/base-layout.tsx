import { Outlet } from "@remix-run/react"
import { Sidebar, SidebarProps } from "~/components/sidebar"

export const BaseLayout = (props: SidebarProps) => {
    return (
        <div className="drawer lg:drawer-open">
            <input id="my-drawer-3" type="checkbox" className="drawer-toggle" />
            <div className="drawer-content flex flex-col">
                <Outlet />
            </div>
            <div className="drawer-side">
                <label htmlFor="my-drawer-3" aria-label="close sidebar" className="drawer-overlay"></label>
                <Sidebar {...props}/>
            </div>
        </div>

    )
}