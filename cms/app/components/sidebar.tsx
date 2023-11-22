import { Link } from "@remix-run/react"
import { IoSettingsOutline, IoDocumentsOutline, IoGlobeOutline } from "react-icons/io5";

export type SidebarProps = {
    items: SidebarItem[]
    siteName: string
}

export type SidebarItem = {
    name: string,
    link: string
}

export const Sidebar = ({items, siteName}: SidebarProps) => {
    return (
        <ul className="menu min-h-full w-80 bg-base-200 p-4">
            <button className="btn btn-ghost mr-auto text-xl">{siteName}</button>
            <div className="flex flex-1 flex-col border-y py-2 my-2">
                {items.map((e) => {
                    return (
                        <li>
                            <Link to={e.link}>
                                <span className="text-xl"><IoDocumentsOutline /></span>
                                {e.name}
                            </Link>
                        </li>
                    )
                })}
            </div>
            <li>
                <a>
                    <span className="text-xl"><IoGlobeOutline /></span> API Playground
                </a>
            </li>
            {/* <li>
                <a><svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg> Commands</a
                >
            </li> */}
            <li>
                <Link to="/settings">
                    <span className="text-xl"><IoSettingsOutline /></span> Site Settings
                </Link>
            </li>
        </ul>
    )
}