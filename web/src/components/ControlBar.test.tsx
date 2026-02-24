import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { ControlBar } from "./ControlBar";

describe("ControlBar", () => {
  it("dispatches step command", async () => {
    const user = userEvent.setup();
    const onCommand = vi.fn();

    render(
      <ControlBar
        policy="rr"
        quantum={2}
        disabled={false}
        onPolicyChange={vi.fn()}
        onQuantumChange={vi.fn()}
        onCommand={onCommand}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Step" }));

    expect(onCommand).toHaveBeenCalledWith({ name: "step", count: 1 });
  });
});
