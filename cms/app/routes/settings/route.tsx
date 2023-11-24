import type { MetaFunction } from "@remix-run/node";
import { Outlet, useLoaderData } from "@remix-run/react";
import { eq } from "drizzle-orm";
import type { Item } from "~/components/selector-bar";
import { SelectorBar } from "~/components/selector-bar";
import { db } from "~/db/database.server";
import { settings } from "~/db/schema";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async () => {
  var values = await db.select().from(settings).all();
  var items: Item[] = values.flatMap((e: any) => { return { id: e.id, name: e.name } })

  const data = { items }
  return data
}

export default function ContentName() {
  const { items } = useLoaderData<any>()
  return (
    <div className="flex flex-row h-full">
      <SelectorBar items={items} title="Settings" linkPath="settings" />
      <Outlet />
    </div>
  )
}