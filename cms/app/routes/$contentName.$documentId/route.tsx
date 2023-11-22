import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { Link, isRouteErrorResponse, useLoaderData, useParams, useRouteError, useSearchParams } from "@remix-run/react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import InputController from "~/components/input-controller";
import Editor from "~/components/richtext/richtext";
import { db } from "~/db/config.server";
import { templates } from "~/utils/templates.server";

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export const loader = async ({ params }: LoaderFunctionArgs) => {
  const newItem = params.id == "new";
  var filteredTemplates = templates.filter(e => e.name == params.contentName?.toLowerCase())
  if (filteredTemplates.length == 0)
  {
    throw new Response(null, {
      status: 404,
      statusText: "Template Not Found",
    });
  }
  const template = filteredTemplates[0];
  console.log(template)

  const data = {
    fields: newItem ? null : await db.prepare(`SELECT * FROM fields where id == ${params.documentId}`).all(),
    newItem,
    template
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
  const { fields, template } = useLoaderData<any>()
  const { handleSubmit, control, formState: { errors } } = useForm<Inputs>({
    mode: 'onChange', values: {
      name: "",
      value: "",
      key: "",
      description: ""
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

            {template.fields.map((e: any) => (
              <InputController
                label={e.name}
                name={e.name}
                control={control}
                defaultValue=""
                rules={{}} />
            ))}
                  <Editor />


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