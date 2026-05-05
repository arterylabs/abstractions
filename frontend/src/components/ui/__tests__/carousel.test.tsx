import { render } from "@testing-library/react";
import { Carousel, CarouselContent, CarouselItem } from "../carousel";

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
});
