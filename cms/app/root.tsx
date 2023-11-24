import { json } from "@remix-run/node";
import CSS from "./css/app.css"
import type { LinksFunction, MetaFunction } from "@remix-run/node";
// import Sidebar from './components/sidebar';
// import { Button, Tooltip } from "flowbite-react";
// import Navbar from "./components/navbar";
import { BaseLayout } from "./layouts/base-layout";
import { templates } from "./utils/templates.server";
import { settings } from "./utils/settings.server";

import {
  Links,
  LiveReload,
  Meta,
  Scripts,
  ScrollRestoration,
  useLoaderData
} from "@remix-run/react";

export const links: LinksFunction = () => {
  return [
    { rel: "stylesheet", href: CSS },
    // { rel: "stylesheet", href: richtext },
  ]
};

export const loader = async () => {
  const links = templates.map(e => ({link: e.name, name: e.displayName}))
  console.log(settings);
  return json({ links, settings });
};

export const meta: MetaFunction = () => {
  return [
    { title: "Dugong CMS" },
    { name: "description", content: "Dugong CMS" },
  ];
};

export default function App() {
  const { links, settings } = useLoaderData<any>();

  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width,initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body className="overscroll-none">
        <BaseLayout items={links} siteName={settings.siteName}/>
        <ScrollRestoration />
        <Scripts />
        {process.env.NODE_ENV === "development" && <LiveReload />}
      </body>
    </html>
  );
}
