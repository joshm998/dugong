import fs from "fs";
import path from "path";
import { singleton } from "./singleton.server";

const getTemplates = () => {
  var templates: any[] = [];
  var templateDir = path.join(process.cwd(), "templates");
  const jsonsInDir = fs.readdirSync(templateDir).filter(file => path.extname(file) === '.json');

  jsonsInDir.forEach(file => {
    const fileData = fs.readFileSync(path.join(templateDir, file));
    const template = JSON.parse(fileData.toString());
    templates.push(template);
  });
  return templates;
}

export const templates = singleton('templates', () => getTemplates())
