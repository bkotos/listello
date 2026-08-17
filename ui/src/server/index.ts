import path from "node:path";
import { fileURLToPath } from "node:url";
import fastifyStatic from "@fastify/static";
import Fastify from "fastify";
import { loadListello } from "@bkotos/listello";

const distDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../dist");

const app = Fastify({ logger: true });

app.get("/health", async () => ({ ok: true }));

app.post<{ Body: { name: string } }>("/api/lists", async (request) => {
  const listello = await loadListello();
  return listello.list.createList(request.body.name);
});

await app.register(fastifyStatic, {
  root: distDir,
});

app.setNotFoundHandler(async (request, reply) => {
  if (request.method !== "GET" || request.url.startsWith("/api/")) {
    return reply.code(404).send({ error: "Not Found" });
  }

  return reply.sendFile("index.html");
});

const port = Number(process.env.PORT ?? 3001);
const host = process.env.HOST ?? "0.0.0.0";

await app.listen({ port, host });
