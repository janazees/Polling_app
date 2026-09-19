import Header from "../components/Header";
import GradientField from "../components/GradientField";
import "./About.css";

export default function About() {
return ( <div className="about-page"> <GradientField /> <Header />

```
  <main className="about-content">

    <section className="about-hero">
      <p className="about-eyebrow">ABOUT POLLS</p>

      <h1>
        Ask a question.
        <br />
        Let the room decide.
      </h1>

      <p className="about-intro">
        Polls is a simple real-time polling platform designed to
        make collecting opinions quick, interactive, and easy to
        follow. Create a question, share one link, and watch the
        results change as people vote.
      </p>
    </section>

    <section className="about-grid">

      <article className="about-card">
        <span className="about-card-number">01</span>

        <h2>Create</h2>

        <p>
          Create a poll with your own question and multiple
          answer options. Your polls are saved to your account
          so you can manage them later.
        </p>
      </article>

      <article className="about-card">
        <span className="about-card-number">02</span>

        <h2>Share</h2>

        <p>
          Share a unique poll link with your audience. Voters
          can join a poll without creating an account.
        </p>
      </article>

      <article className="about-card">
        <span className="about-card-number">03</span>

        <h2>See it live</h2>

        <p>
          Results update in real time as votes arrive, so
          everyone watching the poll can see the latest results
          without refreshing the page.
        </p>
      </article>

    </section>

    <section className="about-tech">

      <div>
        <p className="about-eyebrow">BUILT WITH</p>

        <h2>
          A modern stack
          <br />
          behind a simple experience.
        </h2>
      </div>

      <div className="about-tech-list">
        <span>React</span>
        <span>Go</span>
        <span>Gin</span>
        <span>MongoDB</span>
        <span>Redis</span>
        <span>WebSockets</span>
      </div>

    </section>

    <section className="about-security">

      <div>
        <p className="about-eyebrow">BUILT FOR CONTROL</p>

        <h2>
          Simple to use.
          <br />
          Safe to manage.
        </h2>
      </div>

      <p>
        Poll creation and management are protected by
        authentication, while voting remains accessible through
        a shared poll link. Server-side validation helps keep
        poll data consistent, and voting can be closed by the
        poll creator when the response period ends.
      </p>

    </section>

  </main>
</div>

);
}
