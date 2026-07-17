import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { GoogleOAuthButton } from "../google-oauth-button";

describe("GoogleOAuthButton", () => {
  it("renders the default label", () => {
    render(<GoogleOAuthButton onClick={jest.fn()} />);
    expect(
      screen.getByRole("button", { name: /continue with google/i }),
    ).toBeInTheDocument();
  });

  it("renders a custom label when provided", () => {
    render(
      <GoogleOAuthButton onClick={jest.fn()} label="Sign in with Google" />,
    );
    expect(
      screen.getByRole("button", { name: /sign in with google/i }),
    ).toBeInTheDocument();
  });

  it("calls onClick when the button is clicked", async () => {
    const user = userEvent.setup();
    const onClick = jest.fn();
    render(<GoogleOAuthButton onClick={onClick} />);
    await user.click(
      screen.getByRole("button", { name: /continue with google/i }),
    );
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("renders the Google SVG icon", () => {
    const { container } = render(<GoogleOAuthButton onClick={jest.fn()} />);
    expect(container.querySelector("svg")).toBeInTheDocument();
  });
});
