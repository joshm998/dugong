import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link, isRouteErrorResponse, useLoaderData, useParams, useRouteError, useSearchParams } from "@remix-run/react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import { Input } from "~/components/input";
import InputController from "~/components/input-controller";
import { Item, SelectorBar } from "~/components/selector-bar";
import { db } from "~/db/config.server";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async ({ params }: LoaderFunctionArgs) => {
  console.log(params.id)

  var settings = await db.prepare("SELECT * FROM settings").all();

  var items: Item[] = settings.flatMap((e: any) => { return { id: e.id, name: e.name } })

  const newItem = params.id == "new";

  const data = {
    items,
    document: newItem ? null : await db.prepare(`SELECT * FROM settings where id == ${params.id}`).get(),
    newItem
  }
  return data;
}

type Inputs = {
  name: string
  key: string,
  value: string,
  description: string
}

export default function contentName() {
  const { document, items } = useLoaderData<any>()
  const { handleSubmit, control, formState: { errors } } = useForm<Inputs>({
    mode: 'all', values: {
      name: document.name,
      value: document.value,
      key: document.key,
      description: document.description
    }
  })
  const onSubmit: SubmitHandler<Inputs> = (data) => console.log(data)

  return (
    <>
      <div className="h-full w-full">
        <form onSubmit={handleSubmit((data) => console.log(data))} className="flex gap-2 flex-col w-full h-full p-2 mb-16">
          <div className="flex-1 overflow-scroll max-h-full">
          <InputController
            label="Name"
            name="name"
            control={control}
            defaultValue=""
            rules={{}} />

          <InputController
            label="Key"
            name="key"
            control={control}
            defaultValue=""
            rules={{}} />

          <InputController
            label="Value"
            name="value"
            control={control}
            defaultValue=""
            rules={{}} />
          </div>
          <div className="flex w-full h-16 bg-white p-2 border-t-slate-200 border-t">
            <div className="flex flex-col">
              <p className="my-auto text-xs">Last Published: 22/11/2023</p>
              <p className="my-auto text-xs">Last Updated: 21/11/2023</p>
            </div>

            <input type="submit" value="Save" className="btn ml-auto mr-2" />
            <button className="btn">Publish</button>

          </div>
        </form>
      </div>
    </>
  )
}

export function ErrorBoundary() {
  const error = useRouteError();
  if (isRouteErrorResponse(error)) {
    // error.status = 500
    // error.data = "Oh no! Something went wrong!"
  }
  return (
    <h1>Something went wrong</h1>
  )
}