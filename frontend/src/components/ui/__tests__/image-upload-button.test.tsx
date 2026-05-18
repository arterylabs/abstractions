import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ImageUploadButton } from "../image-upload-button";

describe("ImageUploadButton", () => {
  it("renders 'Upload Image' when no imageUrl is provided", () => {
    render(<ImageUploadButton onFileSelect={jest.fn()} />);
    expect(
      screen.getByRole("button", { name: /upload image/i }),
    ).toBeInTheDocument();
  });

  it("renders 'Change Image' when imageUrl is provided", () => {
    render(
      <ImageUploadButton
        imageUrl="https://example.com/avatar.jpg"
        onFileSelect={jest.fn()}
      />,
    );
    expect(
      screen.getByRole("button", { name: /change image/i }),
    ).toBeInTheDocument();
  });

  it("renders custom upload and change labels", () => {
    const { rerender } = render(
      <ImageUploadButton
        uploadLabel="Add Photo"
        changeLabel="Replace Photo"
        onFileSelect={jest.fn()}
      />,
    );
    expect(
      screen.getByRole("button", { name: /add photo/i }),
    ).toBeInTheDocument();

    rerender(
      <ImageUploadButton
        imageUrl="https://example.com/avatar.jpg"
        uploadLabel="Add Photo"
        changeLabel="Replace Photo"
        onFileSelect={jest.fn()}
      />,
    );
    expect(
      screen.getByRole("button", { name: /replace photo/i }),
    ).toBeInTheDocument();
  });

  it("renders the fallback content when no image is set", () => {
    render(<ImageUploadButton fallback="AB" onFileSelect={jest.fn()} />);
    expect(screen.getByText("AB")).toBeInTheDocument();
  });

  it("calls onFileSelect with the chosen File when a file is selected", () => {
    const onFileSelect = jest.fn();
    render(<ImageUploadButton onFileSelect={onFileSelect} />);
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    const file = new File(["content"], "photo.png", { type: "image/png" });
    fireEvent.change(input, { target: { files: [file] } });
    expect(onFileSelect).toHaveBeenCalledWith(file);
  });

  it("does not call onFileSelect when no file is selected", () => {
    const onFileSelect = jest.fn();
    render(<ImageUploadButton onFileSelect={onFileSelect} />);
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    fireEvent.change(input, { target: { files: [] } });
    expect(onFileSelect).not.toHaveBeenCalled();
  });

  it("shows the loading state on the button when isUploading is true", () => {
    const { container } = render(
      <ImageUploadButton isUploading onFileSelect={jest.fn()} />,
    );
    // Button renders a Spinner (aria-hidden svg) when loading
    expect(
      container.querySelector('svg[aria-hidden="true"]'),
    ).toBeInTheDocument();
  });

  it("clicking the button triggers the file input click", async () => {
    const user = userEvent.setup();
    render(<ImageUploadButton onFileSelect={jest.fn()} />);
    const input = document.querySelector(
      'input[type="file"]',
    ) as HTMLInputElement;
    const clickSpy = jest.spyOn(input, "click");
    await user.click(screen.getByRole("button", { name: /upload image/i }));
    expect(clickSpy).toHaveBeenCalled();
  });
});
