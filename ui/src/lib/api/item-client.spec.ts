import { afterEach, describe, expect, it, vi } from "vitest";
import { addSubtask, commentItem, completeItem, completeSubtask, defineItem, deleteItem, deleteSubtask, getAllItems, getItem, linkAsChild, modifyItemTitle, moveItem, prioritizeSubtask, uncompleteItem, uncompleteSubtask } from "./item-client";
import { request } from "./util";

vi.mock(import("./util"), () => ({
  request: vi.fn(),
}));

afterEach(() => {
  vi.clearAllMocks();
});

describe("defineItem", () => {
  it("posts a new item on a list to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await defineItem("LS_1", { title: "Buy milk" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1/items", {
      method: "POST",
      body: JSON.stringify({ title: "Buy milk" }),
    });
    expect(result).toEqual(item);
  });
});

describe("getItem", () => {
  it("requests an item by id from the API", async () => {
    // Arrange
    const expected = {
      ID: "IT_1",
      ListID: "LS_2",
      ParentID: "",
      Title: "Plan trip",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    const init = { signal: new AbortController().signal };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await getItem("IT_1", init);

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1", init);
    expect(result).toEqual(expected);
  });
});

describe("getAllItems", () => {
  it("requests all items for a list from the API", async () => {
    // Arrange
    const expected = [
      {
        ID: "IT_1",
        ListID: "LS_1",
        ParentID: "",
        Title: "Buy milk",
        Description: "",
        DueDate: "",
        Tags: [],
        Priority: "",
        State: "outstanding",
      },
      {
        ID: "IT_2",
        ListID: "LS_1",
        ParentID: "",
        Title: "Call dentist",
        Description: "",
        DueDate: "",
        Tags: [],
        Priority: "",
        State: "outstanding",
      },
    ];
    const init = { signal: new AbortController().signal };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await getAllItems("LS_1", init);

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1/items", init);
    expect(result).toEqual(expected);
  });
});

describe("completeItem", () => {
  it("posts a complete action for an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "complete",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await completeItem("IT_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/complete", {
      method: "POST",
    });
    expect(result).toEqual(item);
  });
});

describe("uncompleteItem", () => {
  it("posts an uncomplete action for an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await uncompleteItem("IT_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/uncomplete", {
      method: "POST",
    });
    expect(result).toEqual(item);
  });
});

describe("deleteItem", () => {
  it("deletes an item via the API", async () => {
    // Arrange
    vi.mocked(request).mockResolvedValue(undefined);

    // Act
    const result = await deleteItem("IT_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1", {
      method: "DELETE",
    });
    expect(result).toBeUndefined();
  });
});

describe("modifyItemTitle", () => {
  it("patches an item title via the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Schedule dentist",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await modifyItemTitle("IT_1", { title: "Schedule dentist" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/title", {
      method: "PATCH",
      body: JSON.stringify({ title: "Schedule dentist" }),
    });
    expect(result).toEqual(item);
  });
});

describe("moveItem", () => {
  it("posts a move action for an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_2",
      ParentID: "",
      Title: "Schedule dentist",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await moveItem("IT_1", { listID: "LS_2" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/move", {
      method: "POST",
      body: JSON.stringify({ listID: "LS_2" }),
    });
    expect(result).toEqual(item);
  });
});

describe("commentItem", () => {
  it("posts a comment on an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await commentItem("IT_1", { body: "Need 2%" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/comment", {
      method: "POST",
      body: JSON.stringify({ body: "Need 2%" }),
    });
    expect(result).toEqual(item);
  });
});

describe("addSubtask", () => {
  it("posts a subtask on an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await addSubtask("IT_1", { title: "Get 2%" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/subtask", {
      method: "POST",
      body: JSON.stringify({ title: "Get 2%" }),
    });
    expect(result).toEqual(item);
  });
});

describe("completeSubtask", () => {
  it("posts a complete action for a subtask to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await completeSubtask("IT_1", "ST_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/subtask/ST_1/complete", {
      method: "POST",
    });
    expect(result).toEqual(item);
  });
});

describe("uncompleteSubtask", () => {
  it("posts an uncomplete action for a subtask to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await uncompleteSubtask("IT_1", "ST_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/subtask/ST_1/uncomplete", {
      method: "POST",
    });
    expect(result).toEqual(item);
  });
});

describe("deleteSubtask", () => {
  it("deletes a subtask via the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await deleteSubtask("IT_1", "ST_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/subtask/ST_1", {
      method: "DELETE",
    });
    expect(result).toEqual(item);
  });
});

describe("linkAsChild", () => {
  it("posts a link-as-child action for an item to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_2",
      ListID: "LS_1",
      ParentID: "IT_1",
      Title: "Book flights",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await linkAsChild("IT_2", { parentID: "IT_1" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_2/link-as-child", {
      method: "POST",
      body: JSON.stringify({ parentID: "IT_1" }),
    });
    expect(result).toEqual(item);
  });
});

describe("prioritizeSubtask", () => {
  it("posts a prioritize action for a subtask to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await prioritizeSubtask("IT_1", "ST_1", { priority: "high" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/subtask/ST_1/prioritize", {
      method: "POST",
      body: JSON.stringify({ priority: "high" }),
    });
    expect(result).toEqual(item);
  });
});
