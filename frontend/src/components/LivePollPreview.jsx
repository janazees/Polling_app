import { useEffect, useRef, useState } from "react";
import "./LivePollPreview.css";

const QUESTION = "Best time for the team offsite?";
const OPTIONS = ["Weekday morning", "Weekday evening", "Weekend"];

// Deterministic starting split so first paint isn't empty bars.
const INITIAL_VOTES = [18, 9, 13];

/**
 * This is a self-contained visual demo, not a decorative illustration:
 * it shows a viewer exactly what a live poll looks like while votes come
 * in, which is the entire pitch of the product. Once the real /poll/:id
 * page exists, this component's animation logic can be swapped for the
 * actual websocket vote stream.
 */
export default function LivePollPreview() {
  const [votes, setVotes] = useState(INITIAL_VOTES);
  const prefersReducedMotion = useRef(
    typeof window !== "undefined" &&
      window.matchMedia?.("(prefers-reduced-motion: reduce)").matches
  );

  useEffect(() => {
    if (prefersReducedMotion.current) return;

    const interval = setInterval(() => {
      setVotes((current) => {
        const next = [...current];
        const pick = Math.floor(Math.random() * next.length);
        next[pick] += 1;
        return next;
      });
    }, 1700);

    return () => clearInterval(interval);
  }, []);

  const total = votes.reduce((sum, v) => sum + v, 0);

  return (
    <div className="poll-preview" role="img" aria-label={`Live poll preview: ${QUESTION}`}>
      <div className="poll-preview__top">
        <span className="poll-preview__live">
          <span className="poll-preview__dot" />
          Live
        </span>
        <span className="poll-preview__count">{total} votes</span>
      </div>

      <p className="poll-preview__question">{QUESTION}</p>

      <div className="poll-preview__options">
        {OPTIONS.map((label, i) => {
          const pct = total ? Math.round((votes[i] / total) * 100) : 0;
          return (
            <div className="poll-preview__option" key={label}>
              <div className="poll-preview__option-row">
                <span>{label}</span>
                <span className="poll-preview__pct">{pct}%</span>
              </div>
              <div className="poll-preview__bar-track">
                <div
                  className="poll-preview__bar-fill"
                  style={{ width: `${pct}%` }}
                />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
