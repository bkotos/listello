import { loadListello } from "../dist/index.js";

const listello = await loadListello();
const list = listello.list.createList("Next actions");
console.log("created:", list);
