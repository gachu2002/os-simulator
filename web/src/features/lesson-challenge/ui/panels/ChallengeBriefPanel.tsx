interface ChallengeBriefPanelProps {
  objective?: string;
  description?: string;
  successCriteria: string[];
  visualChecks: string[];
}

export function ChallengeBriefPanel({
  objective,
  description,
  successCriteria,
  visualChecks,
}: ChallengeBriefPanelProps) {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-3">
      <h3 className="text-sm font-semibold text-slate-900">Challenge Brief</h3>
      <p className="mt-2 text-sm text-slate-700">
        <span className="font-semibold">Objective:</span> {objective ?? "Understand this lesson behavior."}
      </p>
      <p className="mt-1 text-sm text-slate-600">{description ?? "Run actions and inspect outcomes before submitting."}</p>

      {successCriteria.length > 0 ? (
        <>
          <p className="mt-2 text-xs font-semibold uppercase tracking-wide text-slate-700">Success Criteria</p>
          <ul className="mt-1 grid gap-1">
            {successCriteria.map((item) => (
              <li key={item} className="text-xs text-slate-600">
                - {item}
              </li>
            ))}
          </ul>
        </>
      ) : null}

      {visualChecks.length > 0 ? (
        <>
          <p className="mt-2 text-xs font-semibold uppercase tracking-wide text-slate-700">Visual Checks</p>
          <ul className="mt-1 grid gap-1">
            {visualChecks.map((item) => (
              <li key={item} className="text-xs text-slate-600">
                - {item}
              </li>
            ))}
          </ul>
        </>
      ) : null}
    </section>
  );
}
