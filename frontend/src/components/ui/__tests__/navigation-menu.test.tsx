import { render, screen } from "@testing-library/react";
import {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuTrigger,
  NavigationMenuContent,
  NavigationMenuViewport,
  navigationMenuTriggerStyle,
} from "../navigation-menu";

describe("NavigationMenu", () => {
  it("renders without crashing", () => {
    render(
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <NavigationMenuLink href="#">Home</NavigationMenuLink>
          </NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>,
    );
  });

  it("renders Trigger, Content, and Viewport", () => {
    render(
      <NavigationMenu>
        <NavigationMenuList>
          <NavigationMenuItem>
            <NavigationMenuTrigger>Products</NavigationMenuTrigger>
            <NavigationMenuContent>
              <NavigationMenuLink href="/docs">Docs</NavigationMenuLink>
            </NavigationMenuContent>
          </NavigationMenuItem>
        </NavigationMenuList>
        <NavigationMenuViewport />
      </NavigationMenu>,
    );

    expect(
      screen.getByRole("button", { name: /Products/i }),
    ).toBeInTheDocument();
  });

  it("exports a usable navigationMenuTriggerStyle helper", () => {
    expect(typeof navigationMenuTriggerStyle()).toBe("string");
    expect(navigationMenuTriggerStyle().length).toBeGreaterThan(0);
  });
});
