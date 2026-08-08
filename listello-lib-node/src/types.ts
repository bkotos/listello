/** Mirrors Go domain.List field names as exposed via syscall/js. */
export interface List {
  ID: string;
  Name: string;
}

/** Mirrors Go domain.Item field names as exposed via syscall/js. */
export interface Item {
  ID: string;
  ListID: string;
  ParentID: string;
  Title: string;
  Description: string;
  DueDate: string;
  Tags: string[];
  Priority: string;
  State: string;
}

export interface ListAPI {
  createList(name: string): List;
}

export interface ItemAPI {
  defineItem(listId: string, title: string): Item;
}

export interface Listello {
  list: ListAPI;
  item: ItemAPI;
}

declare global {
  var listello: Listello | undefined;

  class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<void>;
  }
}

export {};
