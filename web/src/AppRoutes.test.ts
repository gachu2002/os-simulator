import { describe, expect, it } from "vitest";

import { parseRoute } from "./lib/routes";

describe("parseRoute", () => {
  it("parses overview route", () => {
    expect(parseRoute("/")).toEqual({ kind: "overview" });
  });

  it("parses learn route with stage query", () => {
    expect(parseRoute("/lesson/l01-sched-rr-basics/learn?stage=2")).toEqual({
      kind: "learn",
      lessonID: "l01-sched-rr-basics",
      stageIndex: 2,
    });
  });

  it("parses challenge route and defaults stage", () => {
    expect(parseRoute("/lesson/l01-sched-rr-basics/challenge")).toEqual({
      kind: "challenge",
      lessonID: "l01-sched-rr-basics",
      stageIndex: 0,
    });
  });

  it("falls back to overview for invalid routes", () => {
    expect(parseRoute("/challenge/l01/0")).toEqual({ kind: "overview" });
  });
});
