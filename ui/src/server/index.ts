import Fastify from "fastify";
import { loadListello } from "@bkotos/listello";

const app = Fastify({ logger: true });

app.get("/health", async () => ({ ok: true }));

app.post<{ Body: { name: string } }>("/api/lists", async (request) => {
  const listello = await loadListello();
  return listello.list.createList(request.body.name);
});

const port = Number(process.env.PORT ?? 3001);
const host = process.env.HOST ?? "0.0.0.0";

await app.listen({ port, host });
