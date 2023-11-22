// import { createRequestHandler } from "@remix-run/express";
// import express from "express";
// import path from "path";
// import * as build from "dugong-cms"

// const app = express();

// // Add the reuqired files for the CMS
// const packageFilePath = new URL(import.meta.resolve('dugong-cms/public')).pathname;
// console.log(packageFilePath)
// app.use(express.static(packageFilePath))
// app.all("*", createRequestHandler({ build }));

// app.listen(3000, () => {
//   console.log("App listening on http://localhost:3000");
// });

const {
  createRequestHandler,
} = require("@remix-run/express");
const express = require("express");

const app = express();

// needs to handle all verbs (GET, POST, etc.)
app.all(
  "*",
  createRequestHandler({
    // `remix build` and `remix dev` output files to a build directory, you need
    // to pass that build to the request handler
    build: require("./build"),

    // return anything you want here to be available as `context` in your
    // loaders and actions. This is where you can bridge the gap between Remix
    // and your server
    getLoadContext(req, res) {
      return {};
    },
  })
);