import { useContext, useState, useEffect } from "react";
import {
    Search,
    Filter,
    Sparkles,
    RefreshCw
} from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import Article from "../components/Article";
import ArticleModal from "../components/ArticleModal";
import api from "../api/api";

export default function DeepSearchPage() {
    const { user } = useContext(AuthContext);
    const isAdmin = user?.role === "admin";

    // Search state
    const [searchText, setSearchText] = useState("");
    const [selectedCategory, setSelectedCategory] = useState("");
    const [selectedSubcategory, setSelectedSubcategory] = useState("");
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
        // load user feedbacks for displayed articles
    }, [articles])

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
                <div className="mb-6 bg-slate-900/60 backdrop-blur-xl rounded-2xl border border-blue-500/20 p-6">
                    <div className="flex flex-col lg:flex-row gap-4">
                        <div className="flex-1 relative">
                            <input
                                type="text"
                                value={searchText}
                                onChange={(e) => setSearchText(e.target.value)}
                                onKeyPress={(e) => e.key === "Enter" && handleSearch()}
                                placeholder="Search articles by title, abstract, or keywords..."
                                className="w-full pl-12 pr-4 py-3 bg-slate-800/50 border border-blue-500/30 rounded-lg text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                            />
                        </div>


                        <button
                            onClick={() => setShowFilters(!showFilters)}
                            className={`flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all ${showFilters
                                    ? "bg-blue-500/20 text-blue-400 border border-blue-500/30"
                                    : "bg-slate-800/50 text-slate-300 border border-blue-500/20 hover:border-blue-500/40"
                                }`}
                        >
                            <Filter className="w-5 h-5" />
                            Filters
                        </button>

                        <button
                            onClick={() => handleSearch(0)}
                            disabled={loading}
                            className="flex items-center gap-2 px-8 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white font-semibold rounded-lg hover:from-blue-600 hover:to-cyan-600 transition-all shadow-lg shadow-blue-500/25 disabled:opacity-50"
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
                        <div className="mt-6 pt-6 border-t border-blue-500/20 grid grid-cols-1 md:grid-cols-2 gap-3">
                            <div>
                                <label className="block text-sm font-medium text-slate-300 mb-2">Category</label>
                                <select
                                    value={selectedCategory}
                                    onChange={(e) => setSelectedCategory(e.target.value)}
                                    className="w-full px-4 py-2 bg-slate-800/50 border border-blue-500/30 rounded-lg text-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
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
                                    <label className="block text-sm font-medium text-slate-300 mb-2">Subcategory</label>
                                    <select
                                        value={selectedSubcategory}
                                        onChange={(e) => setSelectedSubcategory(e.target.value)}
                                        className="w-full px-4 py-2 bg-slate-800/50 border border-blue-500/30 rounded-lg text-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
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
