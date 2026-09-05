import { cleanup, fireEvent, screen, waitFor } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { renderPageWithShellContext } from "../test/renderPageWithShellContext";
import ListPage from "./ListPage";

vi.mock("../lib/api/list-client", () => ({
  getList: vi.fn(),
}));

vi.mock("../lib/api/item-client", () => ({
  getAllItems: vi.fn(),
  defineItem: vi.fn(),
  completeItem: vi.fn(),
  uncompleteItem: vi.fn(),
  deleteItem: vi.fn(),
  modifyItemTitle: vi.fn(),
}));

import { completeItem, defineItem, deleteItem, getAllItems, modifyItemTitle, uncompleteItem } from "../lib/api/item-client";
import { getList } from "../lib/api/list-client";

const sampleItems = [
  {
    ID: "IT_1",
    ListID: "LS_1",
    ParentID: "",
    Title: "Buy windshield wipers for truck",
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
    Title: "Draft weekly status update",
    Description: "",
    DueDate: "",
    Tags: [],
    Priority: "",
    State: "outstanding",
  },
];

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

describe("ListPage", () => {
  it("renders the list title from the API", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Work" })).toBeInTheDocument();
    });
    expect(getList).toHaveBeenCalledWith(
      "LS_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
    expect(screen.getByPlaceholderText("Add to Work…")).toBeInTheDocument();
    expect(screen.getByText("No items yet.")).toBeInTheDocument();
  });

  it("renders a generic title when the list cannot be loaded", async () => {
    // Arrange
    vi.mocked(getList).mockRejectedValue(new Error("not found"));
    vi.mocked(getAllItems).mockResolvedValue([]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_missing",
    });

    // Assert
    await waitFor(() => {
      expect(getList).toHaveBeenCalledWith(
        "LS_missing",
        expect.objectContaining({ signal: expect.any(AbortSignal) }),
      );
    });
    expect(screen.getByRole("heading", { name: "List" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Add to List…")).toBeInTheDocument();
  });

  it("loads items for the list from the API", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledWith(
        "LS_1",
        expect.objectContaining({ signal: expect.any(AbortSignal) }),
      );
    });
  });

  it("renders items from the API", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });

    // Assert
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    expect(screen.getByText("Draft weekly status update")).toBeInTheDocument();
    expect(screen.queryByText("No items yet.")).not.toBeInTheDocument();
  });

  it("renders the item count in the header", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });

    // Assert
    await waitFor(() => {
      expect(screen.getByLabelText("2 items")).toBeInTheDocument();
    });
  });

  it("calls defineItem when Enter is pressed in the capture input", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    vi.mocked(defineItem).mockResolvedValue({
      ID: "IT_3",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    });
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Work" })).toBeInTheDocument();
    });

    // Act
    const input = screen.getByPlaceholderText("Add to Work…");
    fireEvent.change(input, { target: { value: "Buy milk" } });
    fireEvent.keyDown(input, { key: "Enter" });

    // Assert
    expect(defineItem).toHaveBeenCalledWith("LS_1", { title: "Buy milk" });
  });

  it("reloads items after Enter is pressed in the capture input", async () => {
    // Arrange
    const updatedItems = [
      {
        ID: "IT_3",
        ListID: "LS_1",
        ParentID: "",
        Title: "Buy milk",
        Description: "",
        DueDate: "",
        Tags: [],
        Priority: "",
        State: "outstanding",
      },
    ];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce([])
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(defineItem).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Work" })).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    const input = screen.getByPlaceholderText("Add to Work…");
    fireEvent.change(input, { target: { value: "Buy milk" } });
    fireEvent.keyDown(input, { key: "Enter" });

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("calls openSidebar when the open menu button is clicked", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_2", Name: "Personal" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    const { openSidebar } = renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_2",
    });
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Personal" })).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(openSidebar).toHaveBeenCalledOnce();
  });

  it("calls completeItem when the checkbox is clicked", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    vi.mocked(completeItem).mockResolvedValue({
      ...sampleItems[0],
      State: "complete",
    });
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Mark complete" })[0]);

    // Assert
    expect(completeItem).toHaveBeenCalledWith("IT_1");
  });

  it("reloads items after completeItem is called", async () => {
    // Arrange
    const updatedItems = [
      { ...sampleItems[0], State: "complete" },
      sampleItems[1],
    ];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(completeItem).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Mark complete" })[0]);

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("calls uncompleteItem when the checked checkbox is clicked", async () => {
    // Arrange
    const completedItems = [{ ...sampleItems[0], State: "complete" }];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(completedItems);
    vi.mocked(uncompleteItem).mockResolvedValue({
      ...completedItems[0],
      State: "outstanding",
    });
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    expect(uncompleteItem).toHaveBeenCalledWith("IT_1");
  });

  it("reloads items after uncompleteItem is called", async () => {
    // Arrange
    const completedItems = [{ ...sampleItems[0], State: "complete" }];
    const updatedItems = [{ ...sampleItems[0], State: "outstanding" }];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(completedItems)
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(uncompleteItem).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("shows the item as outstanding after completing then uncompleting without a page reload", async () => {
    // Arrange
    const completedItems = [
      { ...sampleItems[0], State: "complete" },
      sampleItems[1],
    ];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(completedItems)
      .mockResolvedValueOnce(sampleItems);
    vi.mocked(completeItem).mockResolvedValue(completedItems[0]);
    vi.mocked(uncompleteItem).mockResolvedValue(sampleItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Mark complete" })[0]);
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    expect(uncompleteItem).toHaveBeenCalledWith("IT_1");
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(3);
    });
    expect(screen.queryByRole("button", { name: "Mark incomplete" })).not.toBeInTheDocument();
    const title = screen.getByText("Buy windshield wipers for truck");
    expect(title).not.toHaveClass("muted");
    expect(title).not.toHaveClass("line-through");
  });

  it("calls deleteItem when Delete is clicked in Task options", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    vi.mocked(deleteItem).mockResolvedValue(undefined);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

    // Assert
    expect(deleteItem).toHaveBeenCalledWith("IT_1");
  });

  it("reloads items after deleteItem is called", async () => {
    // Arrange
    const remainingItems = [sampleItems[1]];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(remainingItems);
    vi.mocked(deleteItem).mockResolvedValue(undefined);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("removes the item from the list after Delete is clicked", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    vi.mocked(deleteItem).mockImplementation(() => new Promise(() => {}));
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

    // Assert
    expect(screen.queryByText("Buy windshield wipers for truck")).not.toBeInTheDocument();
  });

  it("calls modifyItemTitle when Enter is pressed in the rename input", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    vi.mocked(modifyItemTitle).mockResolvedValue({
      ...sampleItems[0],
      Title: "Schedule dentist",
    });
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Rename" }));
    const input = screen.getByDisplayValue("Buy windshield wipers for truck");
    fireEvent.change(input, { target: { value: "Schedule dentist" } });
    fireEvent.keyDown(input, { key: "Enter" });

    // Assert
    expect(modifyItemTitle).toHaveBeenCalledWith("IT_1", { title: "Schedule dentist" });
  });

    it("reloads items after modifyItemTitle is called", async () => {
    // Arrange
    const updatedItems = [{ ...sampleItems[0], Title: "Schedule dentist" }, sampleItems[1]];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(modifyItemTitle).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Rename" }));
    const input = screen.getByDisplayValue("Buy windshield wipers for truck");
    fireEvent.change(input, { target: { value: "Schedule dentist" } });
    fireEvent.keyDown(input, { key: "Enter" });

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("updates the detail panel title after rename from Task options", async () => {
    // Arrange
    const updatedItems = [{ ...sampleItems[0], Title: "Schedule dentist" }, sampleItems[1]];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(modifyItemTitle).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText("Buy windshield wipers for truck"));

    // Act
    fireEvent.click(screen.getAllByRole("button", { name: "Task options" })[0]);
    fireEvent.click(screen.getByRole("menuitem", { name: "Rename" }));
    const input = screen.getAllByDisplayValue("Buy windshield wipers for truck").find(
      (el) => el.tagName === "INPUT",
    )!;
    fireEvent.change(input, { target: { value: "Schedule dentist" } });
    fireEvent.keyDown(input, { key: "Enter" });

    // Assert
    await waitFor(() => {
      const title = screen.getByDisplayValue("Schedule dentist");
      expect(title.tagName).toBe("TEXTAREA");
    });
  });

  it("calls modifyItemTitle when the detail title is blurred", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    vi.mocked(modifyItemTitle).mockResolvedValue({
      ...sampleItems[0],
      Title: "Schedule dentist",
    });
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    fireEvent.click(screen.getByText("Buy windshield wipers for truck"));

    // Act
    const title = screen.getByDisplayValue("Buy windshield wipers for truck");
    fireEvent.change(title, { target: { value: "Schedule dentist" } });
    fireEvent.blur(title);

    // Assert
    expect(modifyItemTitle).toHaveBeenCalledWith("IT_1", { title: "Schedule dentist" });
  });

  it("reloads items after the detail title is blurred", async () => {
    // Arrange
    const updatedItems = [{ ...sampleItems[0], Title: "Schedule dentist" }, sampleItems[1]];
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems)
      .mockResolvedValueOnce(sampleItems)
      .mockResolvedValueOnce(updatedItems);
    vi.mocked(modifyItemTitle).mockResolvedValue(updatedItems[0]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(1);
    });
    fireEvent.click(screen.getByText("Buy windshield wipers for truck"));

    // Act
    const title = screen.getByDisplayValue("Buy windshield wipers for truck");
    fireEvent.change(title, { target: { value: "Schedule dentist" } });
    fireEvent.blur(title);

    // Assert
    await waitFor(() => {
      expect(getAllItems).toHaveBeenCalledTimes(2);
    });
  });

  it("does not show the detail panel before an item is clicked", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue(sampleItems);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });
    await waitFor(() => {
      expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
    });

    // Assert
    expect(document.querySelector("aside.app-detail")).not.toBeInTheDocument();
  });

  describe("when an item is clicked", () => {
    beforeEach(async () => {
      vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
      vi.mocked(getAllItems).mockResolvedValue(sampleItems);
      renderPageWithShellContext(createElement(ListPage), {
        path: "lists/:listId",
        initialEntry: "/lists/LS_1",
      });
      await waitFor(() => {
        expect(screen.getByText("Buy windshield wipers for truck")).toBeInTheDocument();
      });
      fireEvent.click(screen.getByText("Buy windshield wipers for truck"));
    });

    it("opens the detail panel", () => {
      const panel = document.querySelector("aside.app-detail");
      expect(panel).toBeInTheDocument();
      expect(panel).toHaveClass("is-hidden-touch");
    });

    it("marks the clicked row as active", () => {
      const row = screen.getAllByText("Buy windshield wipers for truck")[0].closest(".task-row");
      expect(row).toHaveClass("is-active");
    });

    it("closes when the Close button is clicked", () => {
      // Act
      fireEvent.click(screen.getByRole("button", { name: "Close" }));

      // Assert
      expect(document.querySelector("aside.app-detail")).not.toBeInTheDocument();
    });

    it("unmarks the clicked row as active when the Close button is clicked", () => {
      // Act
      fireEvent.click(screen.getByRole("button", { name: "Close" }));

      // Assert
      const row = screen.getByText("Buy windshield wipers for truck").closest(".task-row");
      expect(row).not.toHaveClass("is-active");
    });

    it("opens the second item in the detail panel when another item is clicked", () => {
      // Act
      fireEvent.click(screen.getByText("Draft weekly status update"));

      // Assert
      expect(screen.getByDisplayValue("Draft weekly status update")).toBeInTheDocument();
    });
  });
});
