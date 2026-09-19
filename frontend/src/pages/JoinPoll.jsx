import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import GradientField from "../components/GradientField";
import "./JoinPoll.css";

export default function JoinPoll() {
  const navigate = useNavigate();

  const [pollInput, setPollInput] = useState("");
  const [error, setError] = useState("");

  function handleJoin(event) {
    event.preventDefault();

    setError("");

    const input = pollInput.trim();

    if (!input) {
      setError("Please enter a poll link or poll code.");
      return;
    }

    let shareCode = input;

    // If the user pasted the full poll URL,
    // extract the share code from it.
    try {
      if (input.startsWith("http://") || input.startsWith("https://")) {
        const url = new URL(input);

        const parts = url.pathname.split("/").filter(Boolean);

        if (
          parts.length >= 2 &&
          parts[parts.length - 2] === "poll"
        ) {
          shareCode = parts[parts.length - 1];
        } else {
          throw new Error("Invalid poll link.");
        }
      }
    } catch {
      setError("Please enter a valid poll link or poll code.");
      return;
    }

    if (!shareCode) {
      setError("Please enter a valid poll code.");
      return;
    }

    navigate(`/poll/${shareCode}`);
  }

  return (
    <div className="join-poll-page">
      <GradientField />
      <Header />

      <main className="join-poll-content">
        <section className="join-poll-card">

          <div className="join-poll-header">
            <p className="join-poll-eyebrow">
              JOIN A POLL
            </p>

            <h1>
              Have your
              <br />
              say.
            </h1>

            <p className="join-poll-description">
              Enter the poll link or code someone shared
              with you. No account needed.
            </p>
          </div>

          <form
            className="join-poll-form"
            onSubmit={handleJoin}
          >
            <div className="join-poll-field">
              <label htmlFor="poll-input">
                Poll link or code
              </label>

              <input
                id="poll-input"
                type="text"
                placeholder="Paste a poll link or enter its code"
                value={pollInput}
                onChange={(event) =>
                  setPollInput(event.target.value)
                }
                autoComplete="off"
              />
            </div>

            {error && (
              <p className="join-poll-error">
                {error}
              </p>
            )}

            <button
              type="submit"
              className="join-poll-button"
            >
              Join poll
            </button>
          </form>

          <p className="join-poll-note">
            You don't need an account to vote.
          </p>

        </section>
      </main>
    </div>
  );
}