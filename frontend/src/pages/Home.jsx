import { Link } from "react-router-dom";
import Header from "../components/Header";
import GradientField from "../components/GradientField";
import LivePollPreview from "../components/LivePollPreview";
import "./Home.css";

export default function Home() {
  return (
    <div className="home">
      <GradientField />
      <Header />

      <main className="home__hero container">
        <div className="home__copy">
          <h1 className="home__heading">
            Let the room decide through polls.
          </h1>
          <p className="home__subheading">
            Create a poll, share one link, and watch votes land in real time —
            no refreshing, no waiting for someone to "recount hands."
          </p>

          <div className="home__actions">
            <Link to="/join" className="btn btn--ghost">
              Join a poll
            </Link>
            <Link to="/create" className="btn btn--primary">
              Create a poll
            </Link>
          </div>

          <p className="home__hint">
            Joining a poll doesn't need an account — creating one does.
          </p>
        </div>

        <div className="home__visual">
          <LivePollPreview />
        </div>
      </main>
    </div>
  );
}
