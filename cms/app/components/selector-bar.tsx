import { Link, useRevalidator } from "@remix-run/react"
export type SelectorBarProps = {
    items: Item[]
    title: string
    linkPath: string
}

export type Item = {
    id: number;
    name: string;
}

export const SelectorBar = ({ items, title, linkPath }: SelectorBarProps) => {
    return (
        <div className="h-full border-r border-r-slate-200">
            <div className="flex px-4 border-b border-b-slate-200">
                <h2 className="text-2xl py-2">{title}</h2>
                <Link to={`/${linkPath}/new`} className="btn btn-sm my-auto ml-auto" replace>
                    Add
                </Link>
            </div>
            <ul className="menu w-56">
                {items.map((e: Item) => (
                    <li key={e.id}>
                        <Link to={`/${linkPath}/${e.id}`} replace>
                            {e.name}
                        </Link>
                    </li>

                ))}

            </ul>
        </div>
    )
}