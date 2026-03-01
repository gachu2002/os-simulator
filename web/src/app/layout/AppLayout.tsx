import { Outlet, useLocation } from "react-router-dom";

export function AppLayout() {
  const location = useLocation();
  const isLessonWorkspaceRoute = location.pathname.includes("/lesson/");

  if (isLessonWorkspaceRoute) {
    return (
      <main className="min-h-screen w-full bg-slate-100 px-3 py-3 md:px-4 md:py-4">
        <Outlet />
      </main>
    );
  }

  return (
    <main className="min-h-screen w-full bg-[radial-gradient(circle_at_top_left,_#e0f2fe_0%,_#f8fafc_40%,_#ffffff_100%)] px-4 py-5 md:px-6">
      <header className="mx-auto max-w-6xl rounded-xl border border-slate-200 bg-white/80 p-5 shadow-sm backdrop-blur">
        <p className="text-xs uppercase tracking-[0.12em] text-sky-700">OSTEP Simulator</p>
        <h1 className="mt-1 text-2xl font-bold text-slate-900 md:text-3xl">Virtualization: CPU Lessons</h1>
        <p className="mt-1 max-w-3xl text-sm text-slate-600 md:text-base">
          Learn core CPU concepts through deterministic interactive challenges built around
          Section 1.
        </p>
      </header>

      <div className="mx-auto mt-4 grid w-full max-w-6xl">
        <Outlet />
      </div>
    </main>
  );
}
