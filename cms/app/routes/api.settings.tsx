import { LoaderFunctionArgs, json } from "@remix-run/node";
import { settings } from "~/utils/settings.server";

export async function loader({
    params,
}: LoaderFunctionArgs) {

    var data = settings;
    return json(data, 200);
}