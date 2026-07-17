import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselPrevious,
  CarouselNext,
} from "../carousel";

describe("Carousel", () => {
  it("renders without crashing", () => {
    render(
      <Carousel>
        <CarouselContent>
          <CarouselItem>Slide 1</CarouselItem>
          <CarouselItem>Slide 2</CarouselItem>
        </CarouselContent>
      </Carousel>,
    );
  });

  it("renders CarouselPrevious and CarouselNext navigation buttons", () => {
    render(
      <Carousel>
        <CarouselContent>
          <CarouselItem>Slide 1</CarouselItem>
          <CarouselItem>Slide 2</CarouselItem>
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>,
    );

    // Icon text ("arrow_back") is included in the accessible name alongside the sr-only label
    expect(
      screen.getByRole("button", { name: /previous slide/i }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /next slide/i }),
    ).toBeInTheDocument();
  });

  it("CarouselPrevious and CarouselNext are disabled when embla reports no scroll available", () => {
    // Without a real embla scroll API in jsdom, canScrollPrev/canScrollNext default to false
    render(
      <Carousel>
        <CarouselContent>
          <CarouselItem>Slide 1</CarouselItem>
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>,
    );

    expect(
      screen.getByRole("button", { name: /previous slide/i }),
    ).toBeDisabled();
    expect(screen.getByRole("button", { name: /next slide/i })).toBeDisabled();
  });

  it("renders with vertical orientation without crashing", () => {
    render(
      <Carousel orientation="vertical">
        <CarouselContent>
          <CarouselItem>Slide 1</CarouselItem>
          <CarouselItem>Slide 2</CarouselItem>
        </CarouselContent>
        <CarouselPrevious />
        <CarouselNext />
      </Carousel>,
    );

    expect(screen.getByRole("region")).toBeInTheDocument();
  });

  it("handles ArrowLeft and ArrowRight keyboard events without crashing", async () => {
    const user = userEvent.setup();
    render(
      <Carousel>
        <CarouselContent>
          <CarouselItem>Slide 1</CarouselItem>
          <CarouselItem>Slide 2</CarouselItem>
        </CarouselContent>
      </Carousel>,
    );

    const region = screen.getByRole("region");
    region.focus();
    await user.keyboard("{ArrowLeft}");
    await user.keyboard("{ArrowRight}");
    // No crash is sufficient — embla is not wired up in jsdom
  });
});
