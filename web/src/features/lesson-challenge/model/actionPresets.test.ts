import { describe, expect, it } from "vitest";

import { toActionPreset } from "./actionPresets";

describe("toActionPreset", () => {
  it("maps execute_instruction to step preset", () => {
    const preset = toActionPreset("execute_instruction");
    expect(preset.label).toBe("Step (1 tick)");
    expect(preset.options?.count).toBe(1);
  });

  it("maps set_policy_fifo_sjf_stcf to fifo policy option", () => {
    const preset = toActionPreset("set_policy_fifo_sjf_stcf");
    expect(preset.label).toBe("Set Policy (FIFO)");
    expect(preset.options?.policy).toBe("fifo");
  });

  it("falls back to raw action label for unknown actions", () => {
    const preset = toActionPreset("unknown_action");
    expect(preset).toEqual({ label: "unknown_action" });
  });
});
