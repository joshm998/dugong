import type { Config } from "drizzle-kit";
 
export default {
  schema: "./app/db/schema.ts",
  out: "../cli/migrations",
} satisfies Config;