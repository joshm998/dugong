import { redirect, type ActionFunctionArgs, type LoaderFunctionArgs, type MetaFunction, json, DataFunctionArgs } from "@remix-run/node";
import { Form, isRouteErrorResponse, useActionData, useLoaderData, useLocation, useRevalidator, useRouteError, useSubmit } from "@remix-run/react";
import { eq } from "drizzle-orm";
import { db } from "~/db/database.server";
import { settings } from "~/utils/settings.server";
import { settings as settingsSchema } from "~/db/schema";
import { z } from "zod";
import { withZod } from "@remix-validated-form/with-zod";
import { ValidatedForm, validationError } from "remix-validated-form";
import InputField from "~/components/input-controller";
import { SubmitButton } from "~/components/submit-button";
import { useEffect } from "react";


export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};


export const validator = withZod(
  z.object({
    name: z
      .string()
      .min(1, { message: "Name is required" }),
    value: z
      .string()
      .min(1, { message: "Value is required" }),
    key: z
      .string()
      .min(1, { message: "Key is required" }),
    description: z
      .string()
  })
);

export const loader = async ({ params }: LoaderFunctionArgs) => {
  console.log(params.id)

  const isNew = params.id == "new";

  let data: any = {
    isNew,
    document: null
  };

  if (!isNew) {
    var documents = await db.select().from(settingsSchema).where(eq(settingsSchema.id, Number(params.id)))
    if (documents.length < 1) {
      console.log("Not Found");
    }

    data.document = documents[0];
  }
  return data;
}

export const action = async ({
  request,
}: DataFunctionArgs) => {
  const result = await validator.validate(
    await request.formData()
  );

  if (result.error) {
    return validationError(result.error);
  }

  type settingNew = typeof settingsSchema.$inferInsert;

  var a: settingNew = {
    name: result.data.name,
    key: result.data.key,
    value: result.data.value,
    description: result.data.description
  }

  var test = await db.insert(settingsSchema).values(a).returning();
  console.log(test);

  redirect(`/settings/${test[0].id}`);
};

type Inputs = {
  name: string
  key: string,
  value: string,
  description: string
}

export default function ContentName() {
  const { document, isNew } = useLoaderData<any>()

  return (
    <>
      <div className="h-full w-full">
        <ValidatedForm
          key={isNew ? "newItem" : document.id}
          validator={validator}
          defaultValues={{
            name: document?.name,
            key: document?.key,
            value: document?.value,
            description: document?.description
          }}
          method="post"
          className="flex gap-2 flex-col w-full h-full p-2 mb-16">
          <div className="flex-1 overflow-scroll max-h-full">
            <InputField name="name" label="Name" />
            <InputField name="key" label="Key" />
            <InputField name="value" label="value" />
            <InputField name="description" label="Description" />
          </div>
          <div className="flex w-full h-16 bg-white p-2 border-t-slate-200 border-t">
            <div className="flex flex-col">
              <p className="my-auto text-xs">Last Published: 22/11/2023</p>
              <p className="my-auto text-xs">Last Updated: 21/11/2023</p>
            </div>

            <input type="submit" value="Save" className="btn ml-auto mr-2" />
            <SubmitButton />
          </div>

        </ValidatedForm>
        {/* <Form method="post" className="flex gap-2 flex-col w-full h-full p-2 mb-16">
          <div className="flex-1 overflow-scroll max-h-full">

            {actionData?.errors.name}

            <Input
              label="Name"
              name="name"
              error={actionData?.errors?.name}
            /> */}

            {/* <InputController
              label="Name"
              name="name"
              control={control}
              defaultValue=""
              rules={{ required: "Required" }} /> */}

            {/* <InputController
              label="Key"
              name="key"
              control={control}
              defaultValue=""
              rules={{ required: "Required" }} />

            <InputController
              label="Value"
              name="value"
              control={control}
              defaultValue=""
              rules={{ required: "Required" }} />

            <InputController
              label="Description"
              name="description"
              control={control}
              defaultValue=""
              rules={{}} /> */}
          {/* </div>
          <div className="flex w-full h-16 bg-white p-2 border-t-slate-200 border-t">
            <div className="flex flex-col">
              <p className="my-auto text-xs">Last Published: 22/11/2023</p>
              <p className="my-auto text-xs">Last Updated: 21/11/2023</p>
            </div>

            <input type="submit" value="Save" className="btn ml-auto mr-2" />
            <button className="btn">Publish</button>

          </div>
        </Form> */}
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