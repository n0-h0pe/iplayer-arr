import { createSignal, Show } from "solid-js";
import { Router, Route } from "@solidjs/router";
import Nav from "./components/Nav";
import Dashboard from "./pages/Dashboard";
import Search from "./pages/Search";
import Config from "./pages/Config";
import Overrides from "./pages/Overrides";

function Layout(props: { children?: any }) {
  return (
    <div class="layout">
      <Nav />
      <main class="main">{props.children}</main>
    </div>
  );
}

function Login(props: { onLogin: () => void }) {
  const [key, setKey] = createSignal("");
  const [error, setError] = createSignal("");

  async function submit(e: Event) {
    e.preventDefault();
    setError("");
    try {
      const res = await fetch("/api/downloads", {
        headers: { "Authorization": `Bearer ${key()}` },
      });
      if (res.ok) {
        localStorage.setItem("iplayer_arr_apikey", key());
        props.onLogin();
      } else {
        setError("Invalid API key");
      }
    } catch {
      setError("Connection failed");
    }
  }

  return (
    <div style="display:flex;align-items:center;justify-content:center;min-height:100vh">
      <div class="card" style="width:360px">
        <h2>iplayer-arr</h2>
        <p style="color:var(--text-muted);margin:0.5rem 0">Enter your API key to continue.</p>
        <form onSubmit={submit}>
          <input
            type="password"
            placeholder="API Key"
            value={key()}
            onInput={e => setKey(e.target.value)}
            style="width:100%;margin:0.5rem 0"
          />
          <Show when={error()}>
            <p style="color:var(--error);font-size:0.85rem">{error()}</p>
          </Show>
          <button type="submit" style="width:100%">Login</button>
        </form>
      </div>
    </div>
  );
}

export default function App() {
  const [authed, setAuthed] = createSignal(!!localStorage.getItem("iplayer_arr_apikey"));

  return (
    <Show when={authed()} fallback={<Login onLogin={() => setAuthed(true)} />}>
      <Router root={Layout}>
        <Route path="/" component={Dashboard} />
        <Route path="/search" component={Search} />
        <Route path="/config" component={Config} />
        <Route path="/overrides" component={Overrides} />
      </Router>
    </Show>
  );
}
