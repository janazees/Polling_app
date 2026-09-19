
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import Header from "../components/Header";
import API_URL from "../api";
import GradientField from "../components/GradientField";
import "./Poll.css";

function getVoterId() {
  let voterId = localStorage.getItem("voterId");

  if (!voterId) {
    voterId = crypto.randomUUID();
    localStorage.setItem("voterId", voterId);
  }

  return voterId;
}

export default function Poll() {
  const { pollId } = useParams();

  const [poll, setPoll] = useState(null);
  const [results, setResults] = useState({});
  const [totalVotes, setTotalVotes] = useState(0);

  const [selectedOption, setSelectedOption] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // --------------------------------------------------
  // LOAD POLL
  // --------------------------------------------------

  useEffect(() => {
    async function fetchPoll() {
      try {
        const response = await fetch(
          `${API_URL}/api/polls/${pollId}`
        );

        const data = await response.json();

        if (!response.ok) {
          throw new Error(
            data.error || "Failed to load poll."
          );
        }

        setPoll(data.poll);
        setResults(data.results || {});
        setTotalVotes(data.totalVotes || 0);
      } catch (err) {
        console.error("Poll loading error:", err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchPoll();
  }, [pollId]);

  // --------------------------------------------------
  // WEBSOCKET - LIVE RESULTS
  // --------------------------------------------------

  useEffect(() => {
    if (!poll) {
      return;
    }

    const socket = new WebSocket(
      `${API_URL.replace(/^http/, "ws")}/api/polls/${pollId}/ws`
    );

    socket.onopen = () => {
      console.log("Live poll connection established");
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.results) {
          setResults(data.results);
        }

        if (typeof data.totalVotes === "number") {
          setTotalVotes(data.totalVotes);
        }
      } catch (err) {
        console.error(
          "Failed to process live update:",
          err
        );
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onclose = () => {
      console.log("Live poll connection closed");
    };

    return () => {
      socket.close();
    };
  }, [poll, pollId]);

  // --------------------------------------------------
  // SUBMIT VOTE
  // --------------------------------------------------

  async function handleVote(event) {
    event.preventDefault();

    if (!selectedOption) {
      setError("Please select an option.");
      return;
    }

    setError("");
    setSuccess("");
    setSubmitting(true);

    try {
      const response = await fetch(
        `${API_URL}/api/polls/${pollId}/vote`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            optionId: selectedOption,
            voterId: getVoterId(),
          }),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        throw new Error(
          data.error || "Failed to submit vote."
        );
      }

      if (data.results) {
        setResults(data.results);
      }

      if (typeof data.totalVotes === "number") {
        setTotalVotes(data.totalVotes);
      }

      setSuccess("Your vote has been recorded!");
      setSelectedOption("");
    } catch (err) {
      console.error("Vote error:", err);
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  // --------------------------------------------------
  // LOADING
  // --------------------------------------------------

  if (loading) {
    return (
      <div className="poll-page">
        <GradientField />
        <Header />

        <main className="poll-content">
          <div className="poll-message">
            Loading poll...
          </div>
        </main>
      </div>
    );
  }

  // --------------------------------------------------
  // ERROR
  // --------------------------------------------------

  if (error && !poll) {
    return (
      <div className="poll-page">
        <GradientField />
        <Header />

        <main className="poll-content">
          <div className="poll-message poll-error">
            {error}
          </div>

          <Link to="/polls" className="poll-back">
            ← Back to polls
          </Link>
        </main>
      </div>
    );
  }

  // --------------------------------------------------
  // POLL PAGE
  // --------------------------------------------------

  return (
    <div className="poll-page">
      <GradientField />
      <Header />

      <main className="poll-content">
        <section className="poll-card-main">

          <p className="poll-eyebrow">
            {poll.isClosed ? "POLL CLOSED" : "LIVE POLL"}
          </p>

          <h1>{poll.question}</h1>

          {poll.isClosed ? (
            <div className="poll-closed">
              <div className="poll-closed-icon">
                ✓
              </div>

              <h2>Voting is closed</h2>

              <p>
                This poll is no longer accepting votes.
                You can still view the final results below.
              </p>
            </div>
          ) : (
            <p className="poll-subtitle">
              Choose one option and submit your vote.
            </p>
          )}

          {/* RESULTS */}

          <div className="poll-options">

            {poll.options.map((option) => {
              const count = results[option.id] || 0;

              const percentage =
                totalVotes > 0
                  ? Math.round(
                      (count / totalVotes) * 100
                    )
                  : 0;

              return poll.isClosed ? (
                <div
                  key={option.id}
                  className="poll-option"
                >
                  <div className="poll-option-content">

                    <div className="poll-option-top">
                      <span className="poll-option-text">
                        {option.text}
                      </span>

                      <span className="poll-option-count">
                        {count}{" "}
                        {count === 1
                          ? "vote"
                          : "votes"}
                      </span>
                    </div>

                    <div className="poll-progress">
                      <div
                        className="poll-progress-fill"
                        style={{
                          width: `${percentage}%`,
                        }}
                      />
                    </div>

                    <span className="poll-percentage">
                      {percentage}%
                    </span>

                  </div>
                </div>
              ) : (
                <label
                  key={option.id}
                  className={`poll-option ${
                    selectedOption === option.id
                      ? "selected"
                      : ""
                  }`}
                >
                  <input
                    type="radio"
                    name="poll-option"
                    value={option.id}
                    checked={
                      selectedOption === option.id
                    }
                    onChange={(event) =>
                      setSelectedOption(
                        event.target.value
                      )
                    }
                  />

                  <div className="poll-option-content">

                    <div className="poll-option-top">
                      <span className="poll-option-text">
                        {option.text}
                      </span>

                      <span className="poll-option-count">
                        {count}{" "}
                        {count === 1
                          ? "vote"
                          : "votes"}
                      </span>
                    </div>

                    <div className="poll-progress">
                      <div
                        className="poll-progress-fill"
                        style={{
                          width: `${percentage}%`,
                        }}
                      />
                    </div>

                    <span className="poll-percentage">
                      {percentage}%
                    </span>

                  </div>
                </label>
              );
            })}

          </div>

          <p className="poll-total">
            {totalVotes}{" "}
            {totalVotes === 1
              ? "total vote"
              : "total votes"}
          </p>

          {/* VOTING FORM */}

          {!poll.isClosed && (
            <form onSubmit={handleVote}>

              {error && (
                <p className="poll-form-error">
                  {error}
                </p>
              )}

              {success && (
                <p className="poll-success">
                  {success}
                </p>
              )}

              <button
                type="submit"
                className="submit-vote-button"
                disabled={submitting}
              >
                {submitting
                  ? "Submitting..."
                  : "Submit vote"}
              </button>

            </form>
          )}

        </section>

        <Link to="/polls" className="poll-back">
          ← Back to your polls
        </Link>
      </main>
    </div>
  );
}
