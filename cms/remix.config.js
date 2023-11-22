/** @type {import('@remix-run/dev').AppConfig} */
export default {
  ignoredRouteFiles: ["**/.*", "dist"],
  // appDirectory: "app",
  assetsBuildDirectory: "public/dist",
  publicPath: "/dist/",
  serverBuildPath: "dist/index.js",
  postcss: true,
  serverDependenciesToBundle: [/^react-icons/],
};
