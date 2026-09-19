import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Header from "../components/Header";
import GradientField from "../components/GradientField";
import "./Polls.css";

export default function Polls() {
  const navigate = useNavigate();

  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadPolls() {
      const token = localStorage.getItem("token");

      if (!token) {
        navigate("/login");
        return;
      }

      try {
        const response = await fetch(
          "http://localhost:8080/api/polls",
          {
            method: "GET",
            headers: {
              Authorization: "Bearer " + token,
            },
          }
        );

        const text = await response.text();

        let data;

        try {
          data = JSON.parse(text);
        } catch {
          throw new Error(
            "Backend returned an invalid response."
          );
        }

        if (!response.ok) {
          throw new Error(
            data.error || "Failed to load polls."
          );
        }

        setPolls(data.polls || []);
      } catch (err) {
        console.error("Poll loading error:", err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    loadPolls();
  }, [navigate]);

  // -----------------------------------------
  // DELETE POLL
  // -----------------------------------------

  async function deletePoll(shareCode) {
    const confirmed = window.confirm(
      "Are you sure you want to delete this poll?"
    );

    if (!confirmed) {
      return;
    }

    const token = localStorage.getItem("token");

    if (!token) {
      navigate("/login");
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/polls/${shareCode}`,
        {
          method: "DELETE",
          headers: {
            Authorization: "Bearer " + token,
          },
        }
      );

      const text = await response.text();

      let data = {};

      if (text) {
        try {
          data = JSON.parse(text);
        } catch {
          throw new Error(
            "Backend returned an invalid response."
          );
        }
      }

      if (!response.ok) {
        throw new Error(
          data.error || "Failed to delete poll."
        );
      }

      // Remove the deleted poll from the screen
      setPolls((currentPolls) =>
        currentPolls.filter(
          (poll) => poll.shareCode !== shareCode
        )
      );

    } catch (err) {
      console.error("Delete poll error:", err);
      alert(err.message);
    }
  }

  // -----------------------------------------
  // COPY POLL LINK
  // -----------------------------------------

  async function copyPollLink(shareCode) {
    const pollLink =
      window.location.origin +
      "/poll/" +
      shareCode;

    try {
      await navigator.clipboard.writeText(pollLink);

      alert("Poll link copied!");
    } catch (err) {
      console.error("Copy error:", err);
      alert("Could not copy the poll link.");
    }
  }

  async function closePoll(shareCode) {
  const confirmed = window.confirm(
    "Are you sure you want to close voting for this poll?"
  );

  if (!confirmed) {
    return;
  }

  const token = localStorage.getItem("token");

  if (!token) {
    navigate("/login");
    return;
  }

  try {
    const response = await fetch(
      `http://localhost:8080/api/polls/${shareCode}/close`,
      {
        method: "POST",
        headers: {
          Authorization: "Bearer " + token,
        },
      }
    );

    const text = await response.text();

    let data = {};

    if (text) {
      try {
        data = JSON.parse(text);
      } catch {
        throw new Error(
          "Backend returned an invalid response."
        );
      }
    }

    if (!response.ok) {
      throw new Error(
        data.error || "Failed to close poll."
      );
    }

    setPolls((currentPolls) =>
      currentPolls.map((poll) =>
        poll.shareCode === shareCode
          ? { ...poll, isClosed: true }
          : poll
      )
    );
  } catch (err) {
    console.error("Close poll error:", err);
    alert(err.message);
  }
}

  return (
    <div className="polls-page">
      <GradientField />
      <Header />

      <main className="polls-content">

        <section className="polls-header">

          <div>
            <p className="polls-eyebrow">
              YOUR POLLS
            </p>

            <h1>
              Your questions,
              <br />
              all in one place.
            </h1>

            <p className="polls-description">
              Manage the polls you've created and share
              them with your audience.
            </p>
          </div>

          <Link
            to="/create"
            className="create-poll-link"
          >
            + Create poll
          </Link>

        </section>

        {loading && (
          <div className="polls-message">
            Loading your polls...
          </div>
        )}

        {error && !loading && (
          <div className="polls-message polls-error">
            {error}
          </div>
        )}

        {!loading && !error && polls.length === 0 && (
          <div className="empty-polls">

            <div className="empty-polls-icon">
              ?
            </div>

            <h2>No polls yet</h2>

            <p>
              You haven't created any polls yet.
              Start with your first question.
            </p>

            <Link
              to="/create"
              className="empty-polls-button"
            >
              Create your first poll
            </Link>

          </div>
        )}

        {!loading && !error && polls.length > 0 && (
          <section className="polls-grid">

            {polls.map((poll) => (
              <article
                className="poll-card"
                key={poll.id}
              >

                <div className="poll-card-top">

                  <span className="poll-card-label">
                    POLL
                  </span>

                  <span className="poll-option-count">
                    {poll.options.length} options
                  </span>

                </div>

                <h2>
                  {poll.question}
                </h2>

                <div className="poll-options-preview">

                  {poll.options
                    .slice(0, 3)
                    .map((option) => (
                      <span key={option.id}>
                        {option.text}
                      </span>
                    ))}

                  {poll.options.length > 3 && (
                    <span>
                      +{poll.options.length - 3} more
                    </span>
                  )}

                </div>

                <div className="poll-card-actions">

                  <Link
                    to={`/poll/${poll.shareCode}`}
                    className="view-poll-button"
                  >
                    {poll.isClosed ? "View results" : "View poll"}
                  </Link>

                  <button
                    type="button"
                    className="copy-poll-button"
                    onClick={() =>
                      copyPollLink(poll.shareCode)
                    }
                  >
                    Copy link
                  </button>

                  {!poll.isClosed && (
                    <button
                      type="button"
                      className="close-poll-button"
                      onClick={() => closePoll(poll.shareCode)}
                    >
                      Close voting
                    </button>
                  )}

                  <button
                    type="button"
                    className="delete-poll-button"
                    onClick={() => deletePoll(poll.shareCode)}
                    aria-label="Delete poll"
                    title="Delete poll"
                  >
                    🗑
                  </button>

                </div>

              </article>
            ))}

          </section>
        )}

      </main>
    </div>
  );
}