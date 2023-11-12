import { createRequestHandler } from "@remix-run/express";
import express from "express";
import path from "path";
import * as build from "dugong-cms"

const app = express();

// Add the reuqired files for the CMS
const packageFilePath = new URL(import.meta.resolve('dugong-cms/public')).pathname;
app.use(express.static(path.dirname(packageFilePath)))
app.all("*", createRequestHandler({ build }));

app.listen(3000, () => {
  console.log("App listening on http://localhost:3000");
});