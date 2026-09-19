import { useState } from "react";
import { Link, NavLink, useNavigate } from "react-router-dom";
import "./Header.css";

const NAV_LINKS = [
  { label: "Home", to: "/" },
  { label: "About", to: "/about" },
  { label: "Browse polls", to: "/polls" },
];

export default function Header() {
  const navigate = useNavigate();

  const [isLoggedIn, setIsLoggedIn] = useState(
    Boolean(localStorage.getItem("token"))
  );

  function handleLogout() {
    localStorage.removeItem("token");
    setIsLoggedIn(false);
    navigate("/login");
  }

  return (
    <header className="site-header">
      <div className="site-header__inner container">
        <Link
          to="/"
          className="site-header__logo"
          aria-label="Polls, go home"
        >
          Polls
        </Link>

        <nav
          className="site-header__nav"
          aria-label="Primary"
        >
          {NAV_LINKS.map((link) => (
            <NavLink
              key={link.to}
              to={link.to}
              className={({ isActive }) =>
                "site-header__link" +
                (isActive ? " is-active" : "")
              }
            >
              {link.label}
            </NavLink>
          ))}
        </nav>

        <div className="site-header__actions">
          {isLoggedIn ? (
            <button
              type="button"
              className="site-header__signin"
              onClick={handleLogout}
            >
              Log out
            </button>
          ) : (
            <Link
              to="/login"
              className="site-header__signin"
            >
              Sign in
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}