// serial runs fn one call at a time: a call made during a run asks for one more run once it ends (several such calls
// make one), so the last state read always wins and no change is dropped. Every caller's promise settles when the
// runs it waited for are over; only an error of the last run rejects it (a failed run with a call waiting runs again).
export function serial(fn: () => Promise<void>): () => Promise<void> {
  let running: Promise<void> | null = null;
  let again = false;
  return () => {
    if (running) { again = true; return running; }
    running = (async () => {
      try {
        for (;;) {
          again = false;
          try { await fn(); } catch (e) { if (!again) throw e; }
          if (!again) return;
        }
      } finally { running = null; again = false; }
    })();
    return running;
  };
}
