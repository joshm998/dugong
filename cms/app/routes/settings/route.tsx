import type { MetaFunction } from "@remix-run/node";
import { Outlet, useLoaderData } from "@remix-run/react";
import { Item, SelectorBar } from "~/components/selector-bar";
import { db } from "~/db/config.server";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async () => {
  var settings = await db.prepare("SELECT * FROM settings").all();

  var items: Item[] = settings.flatMap((e: any) => { return { id: e.id, name: e.name } })

  const data = { items }
  return data
}

export default function contentName() {
  const { items } = useLoaderData<any>()
  return (
    <div className="flex flex-row h-full">
      <SelectorBar items={items} title="Settings" linkPath="settings" />
      <Outlet />
    </div>
  )
}