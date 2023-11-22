import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Outlet, useLoaderData } from "@remix-run/react";
import { SelectorBar } from "~/components/selector-bar";
import { db } from "~/db/config.server";
import { templates } from "~/utils/templates.server";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async ({
  params,
}: LoaderFunctionArgs) => {
  var documents = await db.prepare("SELECT * FROM documents where type == 'project'").all();
  var filteredTemplates = templates.filter(e => e.name == params.contentName?.toLowerCase())
  if (filteredTemplates.length == 0)
  {
    throw new Response(null, {
      status: 404,
      statusText: "Template Not Found",
    });
  }
  const template = filteredTemplates[0];
  var items = documents.map((e: any) => ({id: e.id, name: e.name}))
  const data = {
    items,
    template
  }
  return data
}

export default function contentName() {
  const { items, template } = useLoaderData<any>()

  return (
    <div className="flex flex-row h-full">
      <SelectorBar items={items} title={template.displayName} linkPath={template.name} />
      <Outlet/>
    </div>
  )
}