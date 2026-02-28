import { Outlet } from "react-router-dom";

export function AppLayout() {
  return (
    <main className="mx-auto grid min-h-screen w-full max-w-6xl gap-4 px-4 py-5 md:px-6">
      <header className="space-y-1">
        <p className="text-xs uppercase tracking-[0.12em] text-slate-500">Challenge</p>
        <h1 className="text-2xl font-bold text-slate-900 md:text-3xl">OSTEP Simulator Course</h1>
        <p className="max-w-3xl text-sm text-slate-600 md:text-base">
          Learn operating systems through sectioned lessons and interactive deterministic
          challenges.
        </p>
      </header>

      <Outlet />
    </main>
  );
}
