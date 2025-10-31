import "../assets/styles/global.css";
import { useContext, useState, useEffect } from "react";
import {
  User,
  Mail,
  AtSign,
  Save,
  AlertCircle,
  CheckCircle2,
  Shield,
  Bell,
  Palette,
} from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import api from "../api/api";

const SettingsPage = () => {
  const { user } = useContext(AuthContext);
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [activeTab, setActiveTab] = useState("profile");

  useEffect(() => {
    setProfile(user);
    setLoading(false);
  }, [user]);

  const handleChange = (field, value) => {
    setProfile((prev) => ({ ...prev, [field]: value }));
    setSuccess(false);
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError("");
    setSuccess(false);

    try {
      const updated = {
        first_name: profile.first_name,
        last_name: profile.last_name,
        email: profile.email,
      };
      const data = await api.updateProfile(updated);
      setProfile(data);
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } catch (err) {
      setError(err.message || "Failed to update profile");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div
        className="main-div flex items-center justify-center"
        style={{ minHeight: "60vh" }}
      >
        <div style={{ textAlign: "center", color: "var(--text-secondary)" }}>
          <div
            style={{
              width: "3rem",
              height: "3rem",
              border: "4px solid var(--accent-primary)",
              borderTop: "4px solid transparent",
              borderRadius: "50%",
              margin: "0 auto 1rem",
              animation: "spin 1s linear infinite",
            }}
          />
          <p>Loading your settings...</p>
        </div>
      </div>
    );
  }

  const tabs = [
    { id: "profile", label: "Profile", icon: User },
    { id: "preferences", label: "Preferences", icon: Palette },
    { id: "notifications", label: "Notifications", icon: Bell },
    { id: "security", label: "Security", icon: Shield },
  ];

  return (
    <div
      className="main-div"
      style={{ padding: "2rem 1rem", minHeight: "80vh" }}
    >
      <div className="max-w-5xl mx-auto">
        <div className="mb-4">
          <h1 className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400 mb-2">
            Settings
          </h1>
          <p style={{ color: "var(--text-secondary)" }}>
            Manage your account settings and preferences
          </p>
        </div>

        {error && (
          <div
            className="px-6 py-4 rounded-xl flex items-center gap-3 animate-in fade-in slide-in-from-top-2 mb-4"
            style={{
              background: "rgba(239, 68, 68, 0.1)",
              backdropFilter: "blur(4px)",
              border: "1px solid rgba(239, 68, 68, 0.5)",
              color: "var(--error-text, #f87171)",
            }}
          >
            <AlertCircle size={20} />
            <p>{error}</p>
          </div>
        )}

        {success && (
          <div
            className="mb-4 px-6 py-4 rounded-xl flex items-center gap-3 animate-in fade-in slide-in-from-top-2"
            style={{
              background: "rgba(34, 197, 94, 0.1)",
              backdropFilter: "blur(4px)",
              border: "1px solid rgba(34, 197, 94, 0.5)",
              color: "var(--success-text, #4ade80)",
            }}
          >
            <CheckCircle2 className="w-5 h-5 flex-shrink-0" />
            <p>Profile updated successfully!</p>
          </div>
        )}

        <div style={{ display: "grid", gridTemplateColumns: "1fr 3fr", gap: "1.5rem" }}>
          {/* Sidebar Tabs */}
          <div
            style={{
              background: "var(--card-bg)",
              borderRadius: "16px",
              border: "1px solid var(--border-card)",
              padding: "0.5rem",
              boxShadow: "0 2px 10px var(--shadow-card)",
              display: "flex",
              flexDirection: "column",
              gap: "0.5rem",
            }}
          >
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className="flex items-center gap-3 font-medium transition-all"
                  style={{
                    padding: "0.75rem 1rem",
                    borderRadius: "12px",
                    background: isActive ? "var(--accent-glow)" : "transparent",
                    border: isActive
                      ? "1px solid var(--accent-primary)"
                      : "1px solid transparent",
                    color: isActive
                      ? "var(--accent-primary)"
                      : "var(--text-secondary)",
                    fontWeight: isActive ? 600 : 500,
                  }}
                >
                  <Icon size={18} />
                  {tab.label}
                </button>
              );
            })}
          </div>

          {/* Main Content */}
          <div className="card w-full">
            {activeTab === "profile" && (
              <form
                onSubmit={handleSave}
                className="flex flex-col gap-6"
              >
                <h2 style={{ color: "var(--accent-primary)" }}>
                  Profile Information
                </h2>

                <div>
                  <label
                    className="flex items-center gap-2 text-sm"
                    style={{ color: "var(--text-secondary)" }}
                  >
                    <AtSign size={16} />
                    Username
                  </label>
                  <input
                    type="text"
                    readOnly
                    value={"🔒  " + (profile.username || "")}
                    style={{ opacity: 0.7, cursor: "not-allowed" }}
                  />
                </div>

                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <label
                      className="flex items-center gap-2 text-sm"
                      style={{ color: "var(--text-secondary)" }}
                    >
                      <User size={16} />
                      First Name
                    </label>
                    <input
                      type="text"
                      value={profile.first_name || ""}
                      onChange={(e) =>
                        handleChange("first_name", e.target.value)
                      }
                      placeholder="Enter first name"
                      required
                    />
                  </div>

                  <div>
                    <label
                      className="flex items-center gap-2 text-sm"
                      style={{ color: "var(--text-secondary)" }}
                    >
                      <User size={16} />
                      Last Name
                    </label>
                    <input
                      type="text"
                      value={profile.last_name || ""}
                      onChange={(e) =>
                        handleChange("last_name", e.target.value)
                      }
                      placeholder="Enter last name"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label
                    className="flex items-center gap-2 text-sm"
                    style={{ color: "var(--text-secondary)" }}
                  >
                    <Mail size={16} />
                    Email Address
                  </label>
                  <input
                    type="email"
                    value={profile.email || ""}
                    onChange={(e) => handleChange("email", e.target.value)}
                    placeholder="your.email@example.com"
                    required
                  />
                </div>

                <div
                  className="flex justify-end border-t pt-4"
                  style={{ borderColor: "var(--divider-line)" }}
                >
                  <button
                    type="submit"
                    disabled={saving}
                    className="btn-primary flex items-center gap-2"
                    style={{ width: "fit-content" }}
                  >
                    {saving ? (
                      "Saving..."
                    ) : (
                      <>
                        <Save size={16} /> Save Changes
                      </>
                    )}
                  </button>
                </div>
              </form>
            )}

            {activeTab !== "profile" && (
              <div
                className="text-center py-16"
                style={{ color: "var(--text-muted)" }}
              >
                {activeTab === "preferences" && (
                  <>
                    <Palette size={48} style={{ opacity: 0.5, marginBottom: "1rem" }} />
                    <p>Preferences settings coming soon</p>
                  </>
                )}
                {activeTab === "notifications" && (
                  <>
                    <Bell size={48} style={{ opacity: 0.5, marginBottom: "1rem" }} />
                    <p>Notification settings coming soon</p>
                  </>
                )}
                {activeTab === "security" && (
                  <>
                    <Shield size={48} style={{ opacity: 0.5, marginBottom: "1rem" }} />
                    <p>Security settings coming soon</p>
                  </>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;
