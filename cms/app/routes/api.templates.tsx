import { LoaderFunctionArgs, json } from "@remix-run/node";
import { templates } from "~/utils/templates.server";

export async function loader({
    params,
}: LoaderFunctionArgs) {

    var data = templates;
    return json(data, 200);
}