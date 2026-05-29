import { Children, isValidElement } from "react";
import type { ReactElement, ReactNode } from "react";
import { ImageResponse } from "next/og";
import { GET, runtime } from "./route";

jest.mock("next/og", () => ({
  ImageResponse: jest.fn().mockImplementation(() => {
    return {
      headers: {
        get: (name: string) =>
          name.toLowerCase() === "content-type" ? "image/png" : null,
      },
      text: async () => "generated og image",
    };
  }),
}));

const imageResponseMock = ImageResponse as jest.MockedClass<
  typeof ImageResponse
>;

describe("og-image.png route", () => {
  beforeEach(() => {
    imageResponseMock.mockClear();
  });

  it("runs in the node runtime", () => {
    expect(runtime).toBe("nodejs");
  });

  it("returns a generated 1200x630 image response with the logo embedded", async () => {
    const response = await GET();

    expect(response.headers.get("content-type")).toBe("image/png");
    expect(await response.text()).toBe("generated og image");

    expect(imageResponseMock).toHaveBeenCalledTimes(1);
    const [element, options] = imageResponseMock.mock.calls[0] as [
      ReactElement<{ children: ReactNode }>,
      { width: number; height: number },
    ];
    const logo = Children.toArray(element.props.children).find(
      (child): child is ReactElement<{ alt: string; src: string }> =>
        isValidElement(child) && child.type === "img",
    );

    expect(options).toEqual({ width: 1200, height: 630 });
    expect(logo).toBeDefined();
    expect(logo?.props.alt).toBe("Abstractions");
    expect(logo?.props.src).toMatch(/^data:image\/svg\+xml;base64,/);
  });
});
