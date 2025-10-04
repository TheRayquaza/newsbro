import { useState, useEffect, useContext, useCallback } from "react";
import { useNavigate } from "react-router";
import { FileText, LogOut } from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import api from "../api/api";
const AdminPage = () => {
  const { logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState("users");
  const [users, setUsers] = useState([]);
  const [articles, setArticles] = useState([]);
  const [categories, setCategories] = useState([]);
  const [subcategories, setSubcategories] = useState([]);
  const [feedbackStats, setFeedbackStats] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Article form states
  const [articleTab, setArticleTab] = useState("browse");
  const [selectedCategory, setSelectedCategory] = useState("All");
  const [selectedSubcategory, setSelectedSubcategory] = useState("All");
  const [articleForm, setArticleForm] = useState({
    title: "",
    summary: "",
    content: "",
    category: "",
    subcategory: "",
    link: "",
  });
  const [updateArticleId, setUpdateArticleId] = useState("");
  const [deleteArticleId, setDeleteArticleId] = useState("");

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await api.getUsers();
      setUsers(Array.isArray(data) ? data : []);
    } catch (err) {
      setError("Failed to fetch users: " + err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchArticles = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const params = {};
      if (selectedCategory !== "All") params.category = selectedCategory;
      if (selectedSubcategory !== "All")
        params.subcategory = selectedSubcategory;

      const data = await api.getArticles(params);
      const articlesList = data.data || (Array.isArray(data) ? data : []);
      setArticles(articlesList);
    } catch (err) {
      setError("Failed to fetch articles: " + err.message);
    } finally {
      setLoading(false);
    }
  }, [selectedCategory, selectedSubcategory]);

  const fetchCategories = useCallback(async () => {
    try {
      const data = await api.getCategories();
      const catList = data.data || (Array.isArray(data) ? data : []);
      setCategories(catList.map((c) => c.name || c));
    } catch (err) {
      console.error("Failed to fetch categories:", err);
    }
  }, []);

  const fetchSubcategories = useCallback(async () => {
    try {
      const data = await api.getSubcategories();
      const subcatList = data.data || (Array.isArray(data) ? data : []);
      setSubcategories(subcatList.map((c) => c.name || c));
    } catch (err) {
      console.error("Failed to fetch subcategories:", err);
    }
  }, []);

  const fetchFeedbackStats = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await api.getFeedbackStats();
      setFeedbackStats(data);
    } catch (err) {
      setError("Failed to fetch feedback stats: " + err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (activeTab === "users") fetchUsers();
    else if (activeTab === "articles") {
      fetchArticles();
      fetchCategories();
      fetchSubcategories();
    } else if (activeTab === "feedback") fetchFeedbackStats();
  }, [
    activeTab,
    fetchUsers,
    fetchArticles,
    fetchCategories,
    fetchSubcategories,
    fetchFeedbackStats,
  ]);

  useEffect(() => {
    if (articleTab === "browse") fetchArticles();
  }, [articleTab, fetchArticles]);

  const handleCreateArticle = async () => {
    if (
      !articleForm.title ||
      !articleForm.category ||
      !articleForm.subcategory ||
      !articleForm.link
    ) {
      setError("All fields are required and Link must be a valid URL.");
      return;
    }

    setLoading(true);
    setError("");
    try {
      await api.createArticle(articleForm);
      setArticleForm({
        title: "",
        summary: "",
        content: "",
        category: "",
        subcategory: "",
        link: "",
      });
      alert("✅ Article created successfully!");
      fetchArticles();
    } catch (err) {
      setError("Failed to create article: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateArticle = async () => {
    if (
      !updateArticleId ||
      !articleForm.title ||
      !articleForm.category ||
      !articleForm.subcategory ||
      !articleForm.link
    ) {
      setError("All fields including Article ID are required.");
      return;
    }

    setLoading(true);
    setError("");
    try {
      await api.updateArticle(updateArticleId, articleForm);
      setArticleForm({
        title: "",
        summary: "",
        content: "",
        category: "",
        subcategory: "",
        link: "",
      });
      setUpdateArticleId("");
      alert("✅ Article updated successfully!");
      fetchArticles();
    } catch (err) {
      setError("Failed to update article: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteArticle = async () => {
    if (!deleteArticleId) {
      setError("Please enter an article ID");
      return;
    }

    setLoading(true);
    setError("");
    try {
      await api.deleteArticle(deleteArticleId);
      setDeleteArticleId("");
      alert("✅ Article deleted successfully!");
      fetchArticles();
    } catch (err) {
      setError("Failed to delete article: " + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-black">
      {/* NAVBAR */}
      <nav className="bg-slate-900/80 backdrop-blur-lg border-b border-blue-500/20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-2">
              <FileText className="w-8 h-8 text-blue-400" />
              <span className="text-2xl font-bold text-blue-400">
                NewsBro Admin
              </span>
            </div>
            <button
              onClick={handleLogout}
              className="flex items-center gap-2 bg-red-500/10 text-red-400 px-4 py-2 rounded-lg hover:bg-red-500/20 transition"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>
        </div>
      </nav>

      {/* MAIN DASHBOARD */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-8">
          <h2 className="text-3xl font-bold text-blue-400 mb-6">
            🛠 Admin Dashboard
          </h2>
          <p className="text-slate-300 mb-8">
            Manage users, articles, and feedback
          </p>

          {error && (
            <div className="bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-3 rounded-lg mb-4">
              {error}
            </div>
          )}

          {/* TAB SWITCHER */}
          <div className="flex gap-4 mb-6 border-b border-slate-700">
            {["users", "articles", "feedback"].map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-6 py-3 font-semibold transition ${
                  activeTab === tab
                    ? "text-blue-400 border-b-2 border-blue-400"
                    : "text-slate-400 hover:text-blue-300"
                }`}
              >
                {tab.charAt(0).toUpperCase() + tab.slice(1)}
              </button>
            ))}
          </div>

          {/* USERS TAB */}
          {activeTab === "users" && (
            <div>
              <h3 className="text-xl font-bold text-blue-300 mb-4">👥 Users</h3>
              {loading ? (
                <p className="text-slate-400">Loading users...</p>
              ) : (
                <ul className="space-y-2">
                  {users.map((u) => (
                    <li
                      key={u.id}
                      className="p-3 bg-slate-800/50 rounded-lg border border-slate-700 text-slate-300"
                    >
                      {u.username} — {u.email}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}

          {/* ARTICLES TAB */}
          {activeTab === "articles" && (
            <div>
              {/* SUB-TABS */}
              <div className="flex gap-4 mb-6">
                {["browse", "create", "update", "delete"].map((t) => (
                  <button
                    key={t}
                    onClick={() => setArticleTab(t)}
                    className={`px-4 py-2 rounded-lg transition ${
                      articleTab === t
                        ? "bg-blue-500/20 text-blue-400"
                        : "bg-slate-800/40 text-slate-400 hover:text-blue-300"
                    }`}
                  >
                    {t.charAt(0).toUpperCase() + t.slice(1)}
                  </button>
                ))}
              </div>

              {/* Browse Articles */}
              {articleTab === "browse" && (
                <div>
                  <div className="flex gap-4 mb-4">
                    <select
                      value={selectedCategory}
                      onChange={(e) => setSelectedCategory(e.target.value)}
                      className="px-3 py-2 bg-slate-800 text-slate-300 rounded-lg"
                    >
                      <option value="All">All Categories</option>
                      {categories.map((c) => (
                        <option key={c} value={c}>
                          {c}
                        </option>
                      ))}
                    </select>

                    <select
                      value={selectedSubcategory}
                      onChange={(e) => setSelectedSubcategory(e.target.value)}
                      className="px-3 py-2 bg-slate-800 text-slate-300 rounded-lg"
                    >
                      <option value="All">All Subcategories</option>
                      {subcategories.map((s) => (
                        <option key={s} value={s}>
                          {s}
                        </option>
                      ))}
                    </select>
                  </div>
                  {loading ? (
                    <p className="text-slate-400">Loading articles...</p>
                  ) : (
                    <ul className="space-y-3">
                      {articles.map((a) => (
                        <li
                          key={a.id}
                          className="p-4 bg-slate-800/50 border border-slate-700 rounded-lg text-slate-200"
                        >
                          <h4 className="font-bold text-blue-300">{a.title}</h4>
                          <p>{a.summary}</p>
                          <span className="text-xs text-slate-400">
                            {a.category} → {a.subcategory}
                          </span>
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              )}

              {/* Create Article */}
              {articleTab === "create" && (
                <div className="space-y-4">
                  {[
                    "title",
                    "summary",
                    "content",
                    "category",
                    "subcategory",
                    "link",
                  ].map((field) => (
                    <input
                      key={field}
                      type="text"
                      placeholder={field}
                      value={articleForm[field]}
                      onChange={(e) =>
                        setArticleForm({
                          ...articleForm,
                          [field]: e.target.value,
                        })
                      }
                      className="w-full px-3 py-2 bg-slate-800 text-slate-200 rounded-lg"
                    />
                  ))}
                  <button
                    onClick={handleCreateArticle}
                    className="bg-blue-500/20 text-blue-400 px-4 py-2 rounded-lg hover:bg-blue-500/30"
                  >
                    ➕ Create Article
                  </button>
                </div>
              )}

              {/* Update Article */}
              {articleTab === "update" && (
                <div className="space-y-4">
                  <input
                    type="text"
                    placeholder="Article ID"
                    value={updateArticleId}
                    onChange={(e) => setUpdateArticleId(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-800 text-slate-200 rounded-lg"
                  />
                  {[
                    "title",
                    "summary",
                    "content",
                    "category",
                    "subcategory",
                    "link",
                  ].map((field) => (
                    <input
                      key={field}
                      type="text"
                      placeholder={field}
                      value={articleForm[field]}
                      onChange={(e) =>
                        setArticleForm({
                          ...articleForm,
                          [field]: e.target.value,
                        })
                      }
                      className="w-full px-3 py-2 bg-slate-800 text-slate-200 rounded-lg"
                    />
                  ))}
                  <button
                    onClick={handleUpdateArticle}
                    className="bg-yellow-500/20 text-yellow-400 px-4 py-2 rounded-lg hover:bg-yellow-500/30"
                  >
                    🔄 Update Article
                  </button>
                </div>
              )}

              {/* Delete Article */}
              {articleTab === "delete" && (
                <div className="space-y-4">
                  <input
                    type="text"
                    placeholder="Article ID"
                    value={deleteArticleId}
                    onChange={(e) => setDeleteArticleId(e.target.value)}
                    className="w-full px-3 py-2 bg-slate-800 text-slate-200 rounded-lg"
                  />
                  <button
                    onClick={handleDeleteArticle}
                    className="bg-red-500/20 text-red-400 px-4 py-2 rounded-lg hover:bg-red-500/30"
                  >
                    🗑 Delete Article
                  </button>
                </div>
              )}
            </div>
          )}

          {/* FEEDBACK TAB */}
          {activeTab === "feedback" && (
            <div>
              <h3 className="text-xl font-bold text-blue-300 mb-4">
                📊 Feedback Stats
              </h3>
              {loading ? (
                <p className="text-slate-400">Loading feedback...</p>
              ) : feedbackStats ? (
                <pre className="bg-slate-800/50 p-4 rounded-lg border border-slate-700 text-slate-200">
                  {JSON.stringify(feedbackStats, null, 2)}
                </pre>
              ) : (
                <p className="text-slate-400">No feedback stats available.</p>
              )}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default AdminPage;
