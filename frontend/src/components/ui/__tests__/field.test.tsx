import { render, screen } from "@testing-library/react";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldSet,
  FieldContent,
  FieldTitle,
} from "../field";

describe("Field", () => {
  it("should render with role='group'", () => {
    render(<Field>Content</Field>);
    expect(screen.getByRole("group")).toBeInTheDocument();
  });
});

describe("FieldError", () => {
  it("should render with role='alert' and display the error message", () => {
    render(<FieldError>Something went wrong</FieldError>);
    expect(screen.getByRole("alert")).toHaveTextContent("Something went wrong");
  });

  it("should render multiple errors as a list when errors array has more than 1 item", () => {
    const errors = [{ message: "Error one" }, { message: "Error two" }];
    render(<FieldError errors={errors} />);
    expect(screen.getByRole("list")).toBeInTheDocument();
    expect(screen.getByText("Error one")).toBeInTheDocument();
    expect(screen.getByText("Error two")).toBeInTheDocument();
  });

  it("should render a single error as plain text (no list) when errors has 1 item", () => {
    const errors = [{ message: "Single error" }];
    render(<FieldError errors={errors} />);
    expect(screen.getByRole("alert")).toHaveTextContent("Single error");
    expect(screen.queryByRole("list")).not.toBeInTheDocument();
  });

  it("should return null when no children and no errors", () => {
    const { container } = render(<FieldError />);
    expect(container.firstChild).toBeNull();
  });

  it("should return null when errors array is empty", () => {
    const { container } = render(<FieldError errors={[]} />);
    expect(container.firstChild).toBeNull();
  });
});

describe("FieldDescription", () => {
  it("should render children text", () => {
    render(<FieldDescription>Helper text for the field</FieldDescription>);
    expect(screen.getByText("Helper text for the field")).toBeInTheDocument();
  });
});

describe("FieldSet", () => {
  it("should render as a fieldset element", () => {
    const { container } = render(<FieldSet>Fields</FieldSet>);
    expect(container.querySelector("fieldset")).toBeInTheDocument();
  });
});

describe("FieldGroup", () => {
  it("should render children", () => {
    render(
      <FieldGroup>
        <span>Group content</span>
      </FieldGroup>,
    );
    expect(screen.getByText("Group content")).toBeInTheDocument();
  });
});

describe("FieldContent", () => {
  it("should render children", () => {
    render(
      <FieldContent>
        <span>Field content</span>
      </FieldContent>,
    );
    expect(screen.getByText("Field content")).toBeInTheDocument();
  });
});

describe("FieldTitle", () => {
  it("should render children text", () => {
    render(<FieldTitle>Field Title</FieldTitle>);
    expect(screen.getByText("Field Title")).toBeInTheDocument();
  });
});
