import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { GitHubOAuthButton } from "../github-oauth-button";

describe("GitHubOAuthButton", () => {
  it("renders the default label", () => {
    render(<GitHubOAuthButton onClick={jest.fn()} />);
    expect(
      screen.getByRole("button", { name: /continue with github/i }),
    ).toBeInTheDocument();
  });

  it("renders a custom label when provided", () => {
    render(
      <GitHubOAuthButton onClick={jest.fn()} label="Sign in with GitHub" />,
    );
    expect(
      screen.getByRole("button", { name: /sign in with github/i }),
    ).toBeInTheDocument();
  });

  it("calls onClick when the button is clicked", async () => {
    const user = userEvent.setup();
    const onClick = jest.fn();
    render(<GitHubOAuthButton onClick={onClick} />);
    await user.click(
      screen.getByRole("button", { name: /continue with github/i }),
    );
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("renders the GitHub SVG icon", () => {
    const { container } = render(<GitHubOAuthButton onClick={jest.fn()} />);
    expect(container.querySelector("svg")).toBeInTheDocument();
  });
});
