import { render, screen } from "@testing-library/react";
import { useForm } from "react-hook-form";
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormDescription,
  FormMessage,
} from "../form";

function TestForm() {
  const methods = useForm({ defaultValues: { username: "" } });
  return (
    <Form {...methods}>
      <form>
        <FormField
          control={methods.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <input {...field} aria-label="username" />
              </FormControl>
              <FormDescription>Your public display name.</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}

describe("Form", () => {
  it("renders form with all sub-components without crashing", () => {
    render(<TestForm />);
    expect(screen.getByText("Username")).toBeInTheDocument();
    expect(screen.getByText("Your public display name.")).toBeInTheDocument();
  });
});
