import { render } from "@testing-library/react";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "../input-otp";

describe("InputOTP", () => {
  it("renders without crashing", () => {
    render(
      <InputOTP maxLength={6}>
        <InputOTPGroup>
          <InputOTPSlot index={0} />
          <InputOTPSlot index={1} />
          <InputOTPSlot index={2} />
        </InputOTPGroup>
      </InputOTP>,
    );
  });
});
