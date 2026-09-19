import { Link } from "react-router-dom";
import Header from "../components/Header";
import GradientField from "../components/GradientField";
import "./Placeholder.css";

/**
 * Shared scaffold for pages that aren't built yet (login, signup,
 * create-poll, vote, results...). Swap the children for the real
 * page once the backend endpoints exist — this just keeps routing,
 * header, and background consistent across the app in the meantime.
 */
export default function Placeholder({ title, description }) {
  return (
    <div className="home">
      <GradientField />
      <Header />
      <main className="placeholder container">
        <h1 className="placeholder__title">{title}</h1>
        <p className="placeholder__desc">{description}</p>
        <Link to="/" className="btn btn--ghost">
          Back home
        </Link>
      </main>
    </div>
  );
}
