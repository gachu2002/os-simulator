import { describe, expect, it } from "vitest";

import {
  fromActionCapabilitiesDTO,
  fromActionCapabilityNotesDTO,
} from "./actionCapabilities";

describe("actionCapabilities adapters", () => {
  it("maps action capabilities DTO", () => {
    expect(
      fromActionCapabilitiesDTO({
        supported_now: ["execute_instruction"],
        planned: ["migrate_job"],
      }),
    ).toEqual({
      supportedNow: ["execute_instruction"],
      planned: ["migrate_job"],
    });
  });

  it("maps action capability notes DTO", () => {
    expect(
      fromActionCapabilityNotesDTO({
        execute_instruction: { status: "supported_now", mapped_command: "step" },
        migrate_job: {
          status: "planned",
          reason: "not implemented",
          fallback_action: "step",
        },
      }),
    ).toEqual({
      execute_instruction: { status: "supported_now", mappedCommand: "step" },
      migrate_job: {
        status: "planned",
        reason: "not implemented",
        fallbackAction: "step",
      },
    });
  });

  it("returns undefined for missing DTO inputs", () => {
    expect(fromActionCapabilitiesDTO(undefined)).toBeUndefined();
    expect(fromActionCapabilityNotesDTO(undefined)).toBeUndefined();
  });
});
