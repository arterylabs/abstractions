import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "../tabs";

describe("Tabs", () => {
  it("renders without crashing", () => {
    render(
      <Tabs defaultValue="tab-1">
        <TabsList>
          <TabsTrigger value="tab-1">Tab 1</TabsTrigger>
          <TabsTrigger value="tab-2">Tab 2</TabsTrigger>
        </TabsList>
        <TabsContent value="tab-1">Content 1</TabsContent>
        <TabsContent value="tab-2">Content 2</TabsContent>
      </Tabs>,
    );
  });

  it("renders tab trigger labels", () => {
    render(
      <Tabs defaultValue="tab-1">
        <TabsList>
          <TabsTrigger value="tab-1">First Tab</TabsTrigger>
          <TabsTrigger value="tab-2">Second Tab</TabsTrigger>
        </TabsList>
        <TabsContent value="tab-1">Content 1</TabsContent>
        <TabsContent value="tab-2">Content 2</TabsContent>
      </Tabs>,
    );
    expect(screen.getByRole("tab", { name: "First Tab" })).toBeInTheDocument();
    expect(screen.getByRole("tab", { name: "Second Tab" })).toBeInTheDocument();
  });

  it("shows content for the default selected tab", () => {
    render(
      <Tabs defaultValue="tab-1">
        <TabsList>
          <TabsTrigger value="tab-1">Tab 1</TabsTrigger>
          <TabsTrigger value="tab-2">Tab 2</TabsTrigger>
        </TabsList>
        <TabsContent value="tab-1">Content 1</TabsContent>
        <TabsContent value="tab-2">Content 2</TabsContent>
      </Tabs>,
    );
    expect(screen.getByText("Content 1")).toBeVisible();
  });

  it("switches content when a different tab trigger is clicked", async () => {
    const user = userEvent.setup();
    render(
      <Tabs defaultValue="tab-1">
        <TabsList>
          <TabsTrigger value="tab-1">Tab 1</TabsTrigger>
          <TabsTrigger value="tab-2">Tab 2</TabsTrigger>
        </TabsList>
        <TabsContent value="tab-1">Content 1</TabsContent>
        <TabsContent value="tab-2">Content 2</TabsContent>
      </Tabs>,
    );
    await user.click(screen.getByRole("tab", { name: "Tab 2" }));
    expect(screen.getByText("Content 2")).toBeVisible();
  });

  it("active tab trigger has aria-selected=true", async () => {
    const user = userEvent.setup();
    render(
      <Tabs defaultValue="tab-1">
        <TabsList>
          <TabsTrigger value="tab-1">Tab 1</TabsTrigger>
          <TabsTrigger value="tab-2">Tab 2</TabsTrigger>
        </TabsList>
        <TabsContent value="tab-1">Content 1</TabsContent>
        <TabsContent value="tab-2">Content 2</TabsContent>
      </Tabs>,
    );
    await user.click(screen.getByRole("tab", { name: "Tab 2" }));
    expect(screen.getByRole("tab", { name: "Tab 2" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    expect(screen.getByRole("tab", { name: "Tab 1" })).toHaveAttribute(
      "aria-selected",
      "false",
    );
  });
});
