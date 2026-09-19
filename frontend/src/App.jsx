import { Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Placeholder from "./pages/Placeholder";
import Signup from "./pages/Signup";
import CreatePoll from "./pages/CreatePoll";
import Polls from "./pages/Polls";
import Poll from "./pages/Poll";
import JoinPoll from "./pages/JoinPoll";
import About from "./pages/About";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />

      <Route path="/login" element={<Login />} />

      <Route path="/signup" element={<Signup />} />

      <Route path="/create" element={<CreatePoll />} />

      <Route path="/join" element={<JoinPoll />} />

      <Route path="/poll/:pollId" element={<Poll />} />

      <Route path="/polls" element={<Polls />} />

      <Route path="/about" element={<About />} />
      
    </Routes>
  );
}