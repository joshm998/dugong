import type { MetaFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import { db } from "~/db/config.server";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async () => {
  const data = {
    fields: await db.prepare('SELECT * FROM fields').all()
  }
  return data
}

export default function Index() {
  const { fields } = useLoaderData<any>()

  return (
    <>
      <div>
        <h1>Fields</h1>
      </div>
      <ul>
        {fields.map((field: any) => (
          <li key={field.id}>
            <p>
              <b>{field.name} </b>
               {field.value}</p>
          </li>
        ))}
      </ul>
    </>
  )
}