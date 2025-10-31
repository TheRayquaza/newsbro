import { useContext, useState, useEffect } from "react";
import {
    Search,
    Filter,
    Sparkles,
    RefreshCw,
    Calendar
} from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import Article from "../components/Article";
import ArticleModal from "../components/ArticleModal";
import api from "../api/api";

export default function DeepSearchPage() {
    const { user } = useContext(AuthContext);
    const isAdmin = user?.role === "admin";

    // Initialize dates: begin = today, end = 30 days ago
    const getDefaultBeginDate = () => {
        const today = new Date();
        return today.toISOString().split('T')[0];
    };

    const getDefaultEndDate = () => {
        const date = new Date();
        date.setDate(date.getDate() - 30);
        return date.toISOString().split('T')[0];
    };

    // Search state
    const [searchText, setSearchText] = useState("");
    const [selectedCategory, setSelectedCategory] = useState("");
    const [selectedSubcategory, setSelectedSubcategory] = useState("");
    const [beginDate, setBeginDate] = useState(getDefaultBeginDate());
    const [endDate, setEndDate] = useState(getDefaultEndDate());
    const [limit, _] = useState(20);
    const [offset, setOffset] = useState(0);
    const [totalResults, setTotalResults] = useState(0);

    // Data state
    const [categories, setCategories] = useState([]);
    const [subcategories, setSubcategories] = useState([]);
    const [articles, setArticles] = useState([]);
    const [selectedArticle, setSelectedArticle] = useState(null);
    const [articleFeedback, setArticleFeedback] = useState(null);

    // UI state
    const [loading, setLoading] = useState(false);
    const [showFilters, setShowFilters] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        loadCategories()
    }, []);

    useEffect(() => {
        loadSubcategories(selectedCategory);
        setSelectedSubcategory("");
    }, [selectedCategory]);

    const loadCategories = async () => {
        try {
            const data = await api.getCategories();
            setCategories(data || []);
        } catch (err) {
            console.error("Failed to load categories:", err);
        }
    };

    const loadSubcategories = async (category) => {
        try {
            const data = await api.getSubcategories(category);
            setSubcategories(data || []);
        } catch (err) {
            console.error("Failed to load subcategories:", err);
        }
    };

    const handleSearch = async (newOffset = 0) => {
        setLoading(true);
        setError("");
        try {
            const opts = { limit, offset: newOffset };
            if (searchText) opts.search = searchText;
            if (selectedCategory) opts.category = selectedCategory;
            if (selectedSubcategory) opts.subcategory = selectedSubcategory;
            if (beginDate) opts.beginDate = beginDate;
            if (endDate) opts.endDate = endDate;

            const data = await api.getArticles(opts);
            setArticles(data.articles || []);
            setOffset(newOffset);
            setTotalResults(data.total || 0);
        } catch (err) {
            setError(err.message || "Search failed");
        } finally {
            setLoading(false);
        }
    };

    const handleViewArticle = async (article) => {
        setSelectedArticle(article);
        if (isAdmin) {
            try {
                const feedback = await api.getArticleFeedback(article.id);
                setArticleFeedback(feedback);
            } catch (err) {
                console.error("Failed to load feedback:", err);
            }
        }
    };

    const handleDeleteArticle = async (id) => {
        if (!confirm("Are you sure you want to delete this article?")) return;
        try {
            await api.deleteArticle(id);
            setArticles(articles.filter((a) => a.id !== id));
            setSelectedArticle(null);
        } catch {
            setError("Failed to delete article");
        }
    };

    return (
        <div className="w-full px-4 sm:px-6 lg:px-8 py-8">
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="mb-8">
                    <div className="flex items-center gap-3 mb-2">
                        <Sparkles className="w-8 h-8 text-cyan-400" />
                        <h1 className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400">
                            Search
                        </h1>
                    </div>
                    <p className="text-slate-400">Advanced article search with filters and analytics</p>
                </div>

                {/* Search Bar */}
                <div
                    style={{
                        backgroundColor: "var(--search-container-bg)",
                        borderColor: "var(--nav-border)",
                    }}
                    className="mb-6 backdrop-blur-xl rounded-2xl border p-6"
                >
                    <div className="flex flex-col lg:flex-row gap-4">
                        <div className="flex-1 relative">
                            <input
                                type="text"
                                value={searchText}
                                onChange={(e) => setSearchText(e.target.value)}
                                onKeyPress={(e) => e.key === "Enter" && handleSearch()}
                                placeholder="Search articles by title, abstract, or keywords..."
                                style={{
                                    backgroundColor: "var(--search-input-bg)",
                                    borderColor: "var(--search-input-border)",
                                    color: "var(--search-input-text)",
                                }}
                                className="w-full pl-12 pr-4 py-3 border rounded-lg placeholder:text-[var(--search-input-placeholder)] focus:outline-none focus:ring-2 focus:ring-[var(--search-input-focus-ring)]"
                            />
                        </div>

                        <button
                            onClick={() => setShowFilters(!showFilters)}
                            style={{
                                backgroundColor: showFilters
                                    ? "var(--search-button-active-bg)"
                                    : "var(--search-button-inactive-bg)",
                                color: showFilters
                                    ? "var(--search-button-active-text)"
                                    : "var(--search-button-inactive-text)",
                                borderColor: showFilters
                                    ? "var(--search-button-active-border)"
                                    : "var(--search-button-inactive-border)",
                            }}
                            onMouseEnter={(e) => {
                                if (!showFilters) {
                                    e.currentTarget.style.borderColor = "var(--search-button-active-border)";
                                }
                            }}
                            onMouseLeave={(e) => {
                                if (!showFilters) {
                                    e.currentTarget.style.borderColor = "var(--search-button-inactive-border)";
                                }
                            }}
                            className="flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all border"
                        >
                            <Filter className="w-5 h-5" />
                            Filters
                        </button>

                        <button
                            onClick={() => handleSearch(0)}
                            disabled={loading}
                            style={{
                                background: `linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))`,
                                boxShadow: "0 10px 15px -3px rgba(59, 130, 246, 0.25)",
                            }}
                            onMouseEnter={(e) => {
                                if (!loading) {
                                    e.currentTarget.style.filter = "brightness(1.1)";
                                }
                            }}
                            onMouseLeave={(e) => {
                                e.currentTarget.style.filter = "brightness(1)";
                            }}
                            className="flex items-center gap-2 px-8 py-3 text-white font-semibold rounded-lg transition-all disabled:opacity-50"
                        >
                            {loading ? (
                                <>
                                    <RefreshCw className="w-5 h-5 animate-spin" />
                                    Searching...
                                </>
                            ) : (
                                <>
                                    <Search className="w-5 h-5" />
                                    Search
                                </>
                            )}
                        </button>
                    </div>

                    {showFilters && (
                        <div
                            style={{ borderColor: "var(--nav-border)" }}
                            className="mt-6 pt-6 border-t grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3"
                        >
                            <div>
                                <label
                                    style={{ color: "var(--search-input-text)" }}
                                    className="block text-sm font-medium mb-2"
                                >
                                    Category
                                </label>
                                <select
                                    value={selectedCategory}
                                    onChange={(e) => setSelectedCategory(e.target.value)}
                                    style={{
                                        backgroundColor: "var(--search-input-bg)",
                                        borderColor: "var(--search-input-border)",
                                        color: "var(--search-input-text)",
                                    }}
                                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--search-input-focus-ring)]"
                                >
                                    <option value="">All Categories</option>
                                    {categories.map((cat) => (
                                        <option key={cat} value={cat}>
                                            {cat}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            {selectedCategory && (
                                <div>
                                    <label
                                        style={{ color: "var(--search-input-text)" }}
                                        className="block text-sm font-medium mb-2"
                                    >
                                        Subcategory
                                    </label>
                                    <select
                                        value={selectedSubcategory}
                                        onChange={(e) => setSelectedSubcategory(e.target.value)}
                                        style={{
                                            backgroundColor: "var(--search-input-bg)",
                                            borderColor: "var(--search-input-border)",
                                            color: "var(--search-input-text)",
                                        }}
                                        className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--search-input-focus-ring)]"
                                    >
                                        <option value="">All Subcategories</option>
                                        {subcategories.map((sub) => (
                                            <option key={sub} value={sub}>
                                                {sub}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                            )}

                            <div>
                                <label
                                    style={{ color: "var(--search-input-text)" }}
                                    className="block text-sm font-medium mb-2 flex items-center gap-1"
                                >
                                    <Calendar className="w-4 h-4" />
                                    Begin Date
                                </label>
                                <input
                                    type="date"
                                    value={beginDate}
                                    onChange={(e) => setBeginDate(e.target.value)}
                                    style={{
                                        backgroundColor: "var(--search-input-bg)",
                                        borderColor: "var(--search-input-border)",
                                        color: "var(--search-input-text)",
                                    }}
                                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--search-input-focus-ring)]"
                                />
                            </div>

                            <div>
                                <label
                                    style={{ color: "var(--search-input-text)" }}
                                    className="block text-sm font-medium mb-2 flex items-center gap-1"
                                >
                                    <Calendar className="w-4 h-4" />
                                    End Date
                                </label>
                                <input
                                    type="date"
                                    value={endDate}
                                    onChange={(e) => setEndDate(e.target.value)}
                                    style={{
                                        backgroundColor: "var(--search-input-bg)",
                                        borderColor: "var(--search-input-border)",
                                        color: "var(--search-input-text)",
                                    }}
                                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-[var(--search-input-focus-ring)]"
                                />
                            </div>
                        </div>
                    )}
                </div>

                {error && (
                    <div className="mb-6 bg-red-500/10 border border-red-500/50 text-red-400 px-6 py-4 rounded-xl">
                        {error}
                    </div>
                )}

                {/* Article Modal */}
                {selectedArticle && (
                    <ArticleModal
                        article={selectedArticle}
                        onClose={() => setSelectedArticle(null)}
                        isAdmin={isAdmin}
                        feedback={articleFeedback}
                        onFeedback={async (isPositive) => {
                            try {
                                await api.createArticleFeedback(selectedArticle.id, { value: isPositive });
                                setArticleFeedback(isPositive);
                            } catch (err) {
                                console.error("Failed to submit feedback:", err);
                            }
                        }}
                    />
                )}

                {/* Articles Grid */}
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                    {articles.map((article) => (
                        <Article
                            key={article.id}
                            article={article}
                            onSelect={() => handleViewArticle(article)}
                            isAdmin={isAdmin}
                            onDelete={() => handleDeleteArticle(article.id)}
                        />
                    ))}
                </div>

                {/* Pagination */}
                {articles.length > 0 && (
                    <div className="flex justify-center items-center gap-4 mt-8">
                        <button
                            onClick={() => handleSearch(Math.max(0, offset - limit))}
                            disabled={offset === 0 || loading}
                            className="flex items-center gap-2 px-4 py-2 bg-slate-800/50 border border-blue-500/30 rounded-lg text-slate-300 hover:text-blue-400 disabled:opacity-50"
                        >
                            ← Previous
                        </button>
                        <span className="text-slate-400">
                            Page {Math.floor(offset / limit) + 1}
                            {totalResults ? ` of ${Math.ceil(totalResults / limit)}` : ""}
                        </span>
                        <button
                            onClick={() => handleSearch(offset + limit)}
                            disabled={articles.length < limit || loading}
                            className="flex items-center gap-2 px-4 py-2 bg-slate-800/50 border border-blue-500/30 rounded-lg text-slate-300 hover:text-blue-400 disabled:opacity-50"
                        >
                            Next →
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
