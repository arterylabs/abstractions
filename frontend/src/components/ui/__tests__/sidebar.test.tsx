import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import {
  SidebarProvider,
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuBadge,
  SidebarMenuSkeleton,
  SidebarTrigger,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarSeparator,
  SidebarMenuAction,
} from "../sidebar";

function renderSidebar(props: { defaultOpen?: boolean } = {}) {
  return render(
    <SidebarProvider defaultOpen={props.defaultOpen ?? true}>
      <Sidebar>
        <SidebarHeader>Header</SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton>Dashboard</SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>Footer</SidebarFooter>
      </Sidebar>
      <SidebarTrigger />
    </SidebarProvider>,
  );
}

describe("Sidebar", () => {
  it("renders without crashing", () => {
    render(
      <SidebarProvider>
        <Sidebar>
          <SidebarHeader>Header</SidebarHeader>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton>Menu Item</SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
          <SidebarFooter>Footer</SidebarFooter>
        </Sidebar>
      </SidebarProvider>,
    );
  });

  it("renders header, content, and footer text", () => {
    renderSidebar();
    expect(screen.getByText("Header")).toBeInTheDocument();
    expect(screen.getByText("Footer")).toBeInTheDocument();
    expect(screen.getByText("Dashboard")).toBeInTheDocument();
  });

  it("renders SidebarGroupLabel", () => {
    renderSidebar();
    expect(screen.getByText("Navigation")).toBeInTheDocument();
  });

  it("renders SidebarTrigger button", () => {
    renderSidebar();
    expect(
      screen.getByRole("button", { name: /toggle sidebar/i }),
    ).toBeInTheDocument();
  });

  it("toggles sidebar open/collapsed when SidebarTrigger is clicked", async () => {
    const user = userEvent.setup();
    const { container } = renderSidebar({ defaultOpen: true });

    const wrapper = container.querySelector("[data-slot='sidebar-wrapper']");
    expect(wrapper).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: /toggle sidebar/i }));
    // After toggle the sidebar state should change — the sidebar element gets data-state
    const sidebar = container.querySelector("[data-slot='sidebar']");
    expect(sidebar).toHaveAttribute("data-state", "collapsed");
  });

  it("renders SidebarMenuBadge with text", () => {
    render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton>Inbox</SidebarMenuButton>
                <SidebarMenuBadge>5</SidebarMenuBadge>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    expect(screen.getByText("5")).toBeInTheDocument();
  });

  it("renders SidebarMenuSkeleton without icon", () => {
    const { container } = render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuSkeleton />
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    const skeleton = container.querySelector("[data-sidebar='menu-skeleton']");
    expect(skeleton).toBeInTheDocument();
    expect(
      container.querySelector("[data-sidebar='menu-skeleton-icon']"),
    ).not.toBeInTheDocument();
  });

  it("renders SidebarMenuSkeleton with icon when showIcon=true", () => {
    const { container } = render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuSkeleton showIcon />
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    expect(
      container.querySelector("[data-sidebar='menu-skeleton-icon']"),
    ).toBeInTheDocument();
  });

  it("renders SidebarMenuButton as active", () => {
    const { container } = render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton isActive>Active Item</SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    const button = container.querySelector("[data-slot='sidebar-menu-button']");
    expect(button).toHaveAttribute("data-active", "true");
  });

  it("renders SidebarSeparator", () => {
    const { container } = render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarSeparator />
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    expect(
      container.querySelector("[data-slot='sidebar-separator']"),
    ).toBeInTheDocument();
  });

  it("renders SidebarMenuAction", () => {
    render(
      <SidebarProvider>
        <Sidebar>
          <SidebarContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton>Settings</SidebarMenuButton>
                <SidebarMenuAction aria-label="More options">
                  …
                </SidebarMenuAction>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarContent>
        </Sidebar>
      </SidebarProvider>,
    );
    expect(
      screen.getByRole("button", { name: "More options" }),
    ).toBeInTheDocument();
  });
});
