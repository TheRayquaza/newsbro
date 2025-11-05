import "../assets/styles/global.css";
import forgeIcon from "../assets/icons/forge.png";
import { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { FileText } from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import api from "../api/api";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { user, loading: authLoading, refreshAuth } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!authLoading && user) {
      navigate("/");
    }
  }, [user, authLoading, navigate]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await api.login(email, password);
      refreshAuth();
      navigate("/");
    } catch (err) {
      setError(err.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleForgeOAuth = () => {
    window.open(api.forgeOAuthUrl(), "_blank");
  };

  return (
    <div className="flex justify-center items-center h-screen w-screen main-div">
      <div
        className="card w-full max-w-md"
        style={{ padding: "2rem", maxWidth: "500px" }}
      >
        <div className="card-header text-center mb-8">
          <FileText className="w-10 h-10 text-blue-400 mx-auto mb-2" />
          <h1>NewsBro</h1>
          <p className="subtitle">Welcome back! Please login to continue</p>
        </div>

        <form onSubmit={handleSubmit}>
          {error && (
            <div className="bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-lg mb-4 whitespace-pre-wrap">
              {error}
            </div>
          )}

          <div className="mb-4">
            <div className="w-1/12 mb-2 ml-1 flex justify-center">
              <label>Email</label>
            </div>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              required
              className="w-11/12"
            />
          </div>

          <div className="mb-4">
            <div className="w-1/12 mb-2 ml-1 flex justify-center">
              <label>Password</label>
            </div>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              required
              className="w-11/12"
            />
          </div>

          <button type="submit" disabled={loading} className="btn-primary mb-4">
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <div className="flex justify-end mt-4">
          <a
            href="/forgot-password"
            onClick={(e) => {
              e.preventDefault();
              navigate("/forgot-password");
            }}
            className="text-sm text-blue-500 hover:underline"
          >
            Forgot Password?
          </a>
        </div>

        <div className="divider mb-4">OR</div>

        <button
          onClick={handleForgeOAuth}
          className="btn-secondary mb-4 flex items-center justify-center gap-1"
        >
          <img src={forgeIcon} className="w-9 h-9" alt="Forge Icon" />
          Continue with ForgeID
        </button>

        <div className="back-link">
          <a
            href="/register"
            onClick={(e) => {
              e.preventDefault();
              navigate("/register");
            }}
          >
            Sign Up for an account
          </a>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
