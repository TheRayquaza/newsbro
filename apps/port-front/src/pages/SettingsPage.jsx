import "../assets/styles/global.css";
import { useContext, useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { User, Mail, AtSign, Save, AlertCircle, CheckCircle2, Shield, Bell, Palette } from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import api from "../api/api";

const SettingsPage = () => {
  const { token } = useContext(AuthContext);
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [activeTab, setActiveTab] = useState("profile");
  const navigate = useNavigate();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        setLoading(true);
        const data = await api.getProfile(token);
        setProfile(data);
      } catch (err) {
        setError(err.message || "Failed to load profile");
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [token, navigate]);

  const handleChange = (field, value) => {
    setProfile((prev) => ({ ...prev, [field]: value }));
    setSuccess(false);
  };

  const handleSave = async (e) => {
    e.preventDefault();
    if (!token) return;

    setSaving(true);
    setError("");
    setSuccess(false);
    
    try {
      const updated = {
        first_name: profile.first_name,
        last_name: profile.last_name,
        email: profile.email,
      };
      const data = await api.updateProfile(updated, token);
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
      <div className="w-full min-h-[60vh] flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-blue-400 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-slate-400">Loading your settings...</p>
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
    <div className="w-full px-4 sm:px-6 lg:px-8 py-8">
      <div className="max-w-5xl mx-auto">
        <div className="mb-4">
          <h1 className="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400 mb-2">
            Settings
          </h1>
          <p className="text-slate-400">Manage your account settings and preferences</p>
        </div>

        {error && (
          <div className="mb-4 bg-red-500/10 backdrop-blur-sm border border-red-500/50 text-red-400 px-6 py-4 rounded-xl flex items-center gap-3 animate-in fade-in slide-in-from-top-2">
            <AlertCircle className="w-5 h-5 flex-shrink-0" />
            <p>{error}</p>
          </div>
        )}

        {success && (
          <div className="mb-4 bg-green-500/10 backdrop-blur-sm border border-green-500/50 text-green-400 px-6 py-4 rounded-xl flex items-center gap-3 animate-in fade-in slide-in-from-top-2">
            <CheckCircle2 className="w-5 h-5 flex-shrink-0" />
            <p>Profile updated successfully!</p>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          <div className="lg:col-span-1">
            <div className="bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-2 space-y-1">
              {tabs.map((tab) => {
                const Icon = tab.icon;
                const isActive = activeTab === tab.id;
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${
                      isActive
                        ? "bg-blue-500/20 text-blue-400 border border-blue-500/30"
                        : "text-slate-400 hover:text-slate-300 hover:bg-slate-800/40"
                    }`}
                  >
                    <Icon className="w-5 h-5" />
                    <span className="font-medium">{tab.label}</span>
                  </button>
                );
              })}
            </div>
          </div>

          <div className="lg:col-span-3">
            <div className="bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-8">
              {activeTab === "profile" && (
                <form onSubmit={handleSave} className="space-y-6">
                  <div>
                    <h2 className="text-2xl font-bold text-blue-400 mb-6">Profile Information</h2>

                    <div className="space-y-2 mb-6">
                      <label className="flex items-center gap-2 text-sm font-medium text-slate-300">
                        <AtSign className="w-4 h-4 text-blue-400" />
                        Username
                      </label>
                      <input
                        type="text"
                        readOnly
                        className="w-full bg-slate-800 border border-blue-500/20 rounded-lg px-4 py-3 text-slate-500 cursor-not-allowed opacity-60 select-none"
                        value={"🔒  " + (profile.username || "")}
                        placeholder="username"
                        />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <label className="flex items-center gap-2 text-sm font-medium text-slate-300">
                          <User className="w-4 h-4 text-blue-400" />
                          First Name
                        </label>
                        <input
                          type="text"
                          value={profile.first_name || ""}
                          onChange={(e) => handleChange("first_name", e.target.value)}
                          className="w-full bg-slate-800/50 border border-blue-500/30 rounded-lg px-4 py-3 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-transparent transition"
                          placeholder="Enter first name"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <label className="flex items-center gap-2 text-sm font-medium text-slate-300">
                          <User className="w-4 h-4 text-blue-400" />
                          Last Name
                        </label>
                        <input
                          type="text"
                          value={profile.last_name || ""}
                          onChange={(e) => handleChange("last_name", e.target.value)}
                          className="w-full bg-slate-800/50 border border-blue-500/30 rounded-lg px-4 py-3 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-transparent transition"
                          placeholder="Enter last name"
                          required
                        />
                      </div>
                    </div>

                    <div className="space-y-2 mt-6">
                      <label className="flex items-center gap-2 text-sm font-medium text-slate-300">
                        <Mail className="w-4 h-4 text-blue-400" />
                        Email Address
                      </label>
                      <input
                        type="email"
                        value={profile.email || ""}
                        onChange={(e) => handleChange("email", e.target.value)}
                        className="w-full bg-slate-800/50 border border-blue-500/30 rounded-lg px-4 py-3 text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-transparent transition"
                        placeholder="your.email@example.com"
                        required
                      />
                    </div>
                  </div>

                  {/* Save Button */}
                  <div className="flex justify-end pt-4 border-t border-blue-500/20">
                    <button
                      type="submit"
                      disabled={saving}
                      className={`flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all ${
                        saving
                          ? "bg-slate-700 text-slate-400 cursor-not-allowed"
                          : "bg-gradient-to-r from-blue-500 to-cyan-500 text-white hover:from-blue-600 hover:to-cyan-600 shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40"
                      }`}
                    >
                      {saving ? (
                        <>
                          <div className="w-4 h-4 border-2 border-slate-400 border-t-transparent rounded-full animate-spin"></div>
                          Saving...
                        </>
                      ) : (
                        <>
                          <Save className="w-4 h-4" />
                          Save Changes
                        </>
                      )}
                    </button>
                  </div>
                </form>
              )}

              {activeTab === "preferences" && (
                <div className="space-y-6">
                  <h2 className="text-2xl font-bold text-blue-400 mb-6">Preferences</h2>
                  <div className="text-center py-12 text-slate-400">
                    <Palette className="w-16 h-16 mx-auto mb-4 opacity-50" />
                    <p className="text-lg">Preferences settings coming soon</p>
                  </div>
                </div>
              )}

              {activeTab === "notifications" && (
                <div className="space-y-6">
                  <h2 className="text-2xl font-bold text-blue-400 mb-6">Notifications</h2>
                  <div className="text-center py-12 text-slate-400">
                    <Bell className="w-16 h-16 mx-auto mb-4 opacity-50" />
                    <p className="text-lg">Notification settings coming soon</p>
                  </div>
                </div>
              )}

              {activeTab === "security" && (
                <div className="space-y-6">
                  <h2 className="text-2xl font-bold text-blue-400 mb-6">Security</h2>
                  <div className="text-center py-12 text-slate-400">
                    <Shield className="w-16 h-16 mx-auto mb-4 opacity-50" />
                    <p className="text-lg">Security settings coming soon</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;