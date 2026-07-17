import { render, screen } from "@testing-library/react";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "../card";

describe("Card", () => {
  it("renders without crashing", () => {
    render(
      <Card>
        <CardHeader>
          <CardTitle>Title</CardTitle>
          <CardDescription>Description</CardDescription>
        </CardHeader>
        <CardContent>Content</CardContent>
        <CardFooter>Footer</CardFooter>
      </Card>,
    );
  });

  it("CardTitle renders its text", () => {
    render(<CardTitle>My Card Title</CardTitle>);
    expect(screen.getByText("My Card Title")).toBeInTheDocument();
  });

  it("CardDescription renders its text", () => {
    render(<CardDescription>A short description</CardDescription>);
    expect(screen.getByText("A short description")).toBeInTheDocument();
  });

  it("CardHeader renders its children", () => {
    render(
      <CardHeader>
        <span>Header child</span>
      </CardHeader>,
    );
    expect(screen.getByText("Header child")).toBeInTheDocument();
  });

  it("CardContent renders its children", () => {
    render(
      <CardContent>
        <p>Body content</p>
      </CardContent>,
    );
    expect(screen.getByText("Body content")).toBeInTheDocument();
  });

  it("CardFooter renders its children", () => {
    render(
      <CardFooter>
        <button>Save</button>
      </CardFooter>,
    );
    expect(screen.getByRole("button", { name: "Save" })).toBeInTheDocument();
  });
});
