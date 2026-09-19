import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import Header from "../components/Header";
import GradientField from "../components/GradientField";
import API_URL from "../api";
import "./CreatePoll.css";

export default function CreatePoll() {
  const navigate = useNavigate();

  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);

  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const [createdPoll, setCreatedPoll] = useState(null);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("token");

    if (!token) {
      navigate("/login");
    }
  }, [navigate]);

  function handleOptionChange(index, value) {
    const updatedOptions = [...options];
    updatedOptions[index] = value;
    setOptions(updatedOptions);
  }

  function addOption() {
    if (options.length < 10) {
      setOptions([...options, ""]);
    }
  }

  function removeOption(index) {
    if (options.length <= 2) {
      return;
    }

    setOptions(
      options.filter(
        (_, optionIndex) => optionIndex !== index
      )
    );
  }

  async function handleCreatePoll(event) {
    event.preventDefault();

    setError("");

    const token = localStorage.getItem("token");

    if (!token) {
      navigate("/login");
      return;
    }

    setLoading(true);

    try {
      const response = await fetch(
        `${API_URL}/api/polls`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            question,
            options: options.map((option) =>
              option.trim()
            ),
          }),
        }
      );

      const data = await response.json();

      if (!response.ok) {
        console.log(
          "Create poll backend response:",
          data
        );

        throw new Error(
          data.error || "Failed to create poll"
        );
      }

      console.log("Poll created:", data);

      setCreatedPoll(data.poll);

    } catch (err) {
      console.error("Create poll error:", err);

      setError(
        err.message || "Something went wrong."
      );
    } finally {
      setLoading(false);
    }
  }

  async function copyPollLink() {
    if (!createdPoll) {
      return;
    }

    const pollLink =
      `${window.location.origin}/poll/` +
      createdPoll.shareCode;

    try {
      await navigator.clipboard.writeText(pollLink);

      setCopied(true);

      setTimeout(() => {
        setCopied(false);
      }, 2000);

    } catch (err) {
      console.error("Copy link error:", err);
      setError("Could not copy the poll link.");
    }
  }

  function openPoll() {
    if (!createdPoll) {
      return;
    }

    navigate(`/poll/${createdPoll.shareCode}`);
  }

  if (createdPoll) {
    const pollLink =
      `${window.location.origin}/poll/` +
      createdPoll.shareCode;

    return (
      <div className="create-poll-page">
        <GradientField />
        <Header />

        <main className="create-poll-content">
          <section className="create-poll-card share-poll-card">

            <div className="share-success-icon">
              ✓
            </div>

            <div className="create-poll-header share-poll-header">
              <p className="create-poll-eyebrow">
                POLL CREATED
              </p>

              <h1>
                Your poll is
                <br />
                ready to share.
              </h1>

              <p className="create-poll-description">
                Send this link to anyone you want to
                invite to vote. They don't need an account.
              </p>
            </div>

            <div className="share-question">
              <span>Your question</span>

              <p>
                {createdPoll.question}
              </p>
            </div>

            <div className="share-link-section">
              <label htmlFor="poll-share-link">
                Share link
              </label>

              <div className="share-link-row">
                <input
                  id="poll-share-link"
                  type="text"
                  value={pollLink}
                  readOnly
                />

                <button
                  type="button"
                  onClick={copyPollLink}
                  className="copy-share-button"
                >
                  {copied ? "Copied!" : "Copy"}
                </button>
              </div>
            </div>

            <button
              type="button"
              className="create-poll-button"
              onClick={openPoll}
            >
              Open poll
            </button>

            <Link
              to="/polls"
              className="create-poll-back"
            >
              ← Back to my polls
            </Link>

          </section>
        </main>
      </div>
    );
  }

  return (
    <div className="create-poll-page">
      <GradientField />
      <Header />

      <main className="create-poll-content">
        <section className="create-poll-card">

          <div className="create-poll-header">
            <p className="create-poll-eyebrow">
              CREATE A POLL
            </p>

            <h1>
              Ask something
              <br />
              worth answering.
            </h1>

            <p className="create-poll-description">
              Create a question, add your options, and share it
              with your audience.
            </p>
          </div>

          <form
            className="create-poll-form"
            onSubmit={handleCreatePoll}
          >

            <div className="create-poll-field">
              <label htmlFor="question">
                Question
              </label>

              <textarea
                id="question"
                placeholder="What do you want to ask?"
                value={question}
                onChange={(event) =>
                  setQuestion(event.target.value)
                }
                maxLength={200}
                required
              />

              <span className="character-count">
                {question.length}/200
              </span>
            </div>

            <div className="options-section">

              <div className="options-header">
                <label>Options</label>

                <span>
                  {options.length}/10
                </span>
              </div>

              <div className="options-list">

                {options.map((option, index) => (
                  <div
                    className="option-input-wrapper"
                    key={index}
                  >
                    <input
                      type="text"
                      placeholder={`Option ${index + 1}`}
                      value={option}
                      onChange={(event) =>
                        handleOptionChange(
                          index,
                          event.target.value
                        )
                      }
                      maxLength={100}
                      required
                    />

                    {options.length > 2 && (
                      <button
                        type="button"
                        className="remove-option"
                        onClick={() =>
                          removeOption(index)
                        }
                        aria-label={`Remove option ${
                          index + 1
                        }`}
                      >
                        ×
                      </button>
                    )}
                  </div>
                ))}

              </div>

              {options.length < 10 && (
                <button
                  type="button"
                  className="add-option"
                  onClick={addOption}
                >
                  <span>+</span>
                  Add option
                </button>
              )}

            </div>

            {error && (
              <p className="create-poll-error">
                {error}
              </p>
            )}

            <button
              type="submit"
              className="create-poll-button"
              disabled={loading}
            >
              {loading
                ? "Creating poll..."
                : "Create poll"}
            </button>

          </form>

          <Link
            to="/"
            className="create-poll-back"
          >
            ← Back home
          </Link>

        </section>
      </main>
    </div>
  );
}
