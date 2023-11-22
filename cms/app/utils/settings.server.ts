import fs from "fs";
import path from "path";
import { singleton } from "./singleton.server";

const getSettings = () => {
  var settingsFile = path.join(process.cwd(), "settings.json");
  const settingsFileData = fs.readFileSync(settingsFile)
  return JSON.parse(settingsFileData.toString());
}

export const settings = singleton('settings', () => getSettings())
