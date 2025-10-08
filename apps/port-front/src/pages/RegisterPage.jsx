import "../assets/styles/global.css";
import React, { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { FileText } from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import api from "../api/api";

const RegisterPage = () => {
  const [formData, setFormData] = useState({
    email: "",
    firstName: "",
    lastName: "",
    username: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { user, loading: authLoading } = useContext(AuthContext);
  let navigate = useNavigate();

  useEffect(() => {
    if (!authLoading && user) {
      navigate("/");
    }
  }, [user, authLoading, navigate]);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await api.register({
        email: formData.email,
        first_name: formData.firstName,
        last_name: formData.lastName,
        username: formData.username,
        password: formData.password,
      });
      navigate("/");
    } catch (err) {
      setError(err.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        width: "100vw",
      }}
      className="main-div"
    >
      <div
        className="card w-full mx-8"
        style={{ padding: "2rem", maxWidth: "700px" }}
      >
        <div className="card-header text-center mb-8">
          <FileText className="w-10 h-10 text-blue-400 mx-auto mb-2" />
          <h1>NewsBro</h1>
          <p className="subtitle">Create your account to get started</p>
        </div>

        <form onSubmit={handleSubmit}>
          {error && (
            <div className="bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-lg mb-4">
              {error}
            </div>
          )}

          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <div className="w-full mb-2 ml-1 justify-start flex">
                <label>First Name</label>
              </div>
              <input
                type="text"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                placeholder="Enter your first name"
                required
                style={{ width: "100%" }}
              />
            </div>

            <div>
              <div className="w-full mb-2 ml-1 justify-start flex">
                <label>Last Name</label>
              </div>
              <input
                type="text"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                placeholder="Enter your last name"
                required
                style={{ width: "100%" }}
              />
            </div>
          </div>

          <div className="mb-4">
            <div className="w-1/12 mb-2 flex justify-center">
              <label>Email</label>
            </div>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="Enter your email"
              required
              style={{ width: "100%" }}
            />
          </div>

          <div>
            <div className="w-1/12 mb-2 ml-5 flex justify-center">
              <label>Username</label>
            </div>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="Choose a username"
              required
              style={{ width: "100%" }}
            />
          </div>

          <div>
            <div className="w-1/12 mb-2 ml-4 flex justify-center">
              <label>Password</label>
            </div>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="Create a password"
              required
              style={{ width: "100%" }}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="btn-primary mb-4 mt-4"
          >
            {loading ? "Creating account..." : "Register"}
          </button>
        </form>

        <div className="back-link">
          <a
            href="/"
            onClick={(e) => {
              e.preventDefault();
              navigate("/");
            }}
          >
            ↩ Back to Login
          </a>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
