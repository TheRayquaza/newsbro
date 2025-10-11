import { useContext, useState, useEffect, useRef } from "react";
import {
    History,
    RefreshCw,
    ChevronDown
} from "lucide-react";
import { AuthContext } from "../contexts/Auth";
import Article from "../components/Article";
import ArticleModal from "../components/ArticleModal";
import api from "../api/api";

const LIMIT = 12;

const HistoryPage = () => {
    const { user } = useContext(AuthContext);
    const isAdmin = user?.role === "admin";

    // Data state
    const [articles, setArticles] = useState([]);
    const [selectedArticle, setSelectedArticle] = useState(null);
    const [offset, setOffset] = useState(0);
    const [totalResults, setTotalResults] = useState(0);

    // UI state
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [error, setError] = useState("");
    const hasMore = articles.length < totalResults;
    const containerRef = useRef(null);

    // Initial load
    useEffect(() => {
        loadHistory(0);
    }, []);

    const loadHistory = async (newOffset) => {
        try {
            const opts = { limit: LIMIT, offset: newOffset };
            const data = await api.getArticleHistory(opts);
            
            if (newOffset === 0) {
                setArticles(data.articles || []);
            } else {
                setArticles(prev => [...prev, ...(data.articles || [])]);
            }
            
            setOffset(newOffset);
            setTotalResults(data.total || 0);
            setError("");
        } catch (err) {
            setError(err.message || "Failed to load history");
        } finally {
            setLoading(false);
            setLoadingMore(false);
        }
    };

    const handleLoadMore = () => {
        setLoadingMore(true);
        loadHistory(offset + LIMIT);
    };

    const handleSelectArticle = (article) => {
        setSelectedArticle(article);
    };

    const handleCloseModal = () => {
        setSelectedArticle(null);
    };

    if (loading) {
        return (
            <div className="w-full px-4 sm:px-6 lg:px-8 py-8 flex items-center justify-center min-h-screen">
                <div className="flex flex-col items-center gap-3">
                    <RefreshCw className="w-8 h-8 text-blue-400 animate-spin" />
                    <p className="text-slate-400">Loading your history...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full px-4 sm:px-6 lg:px-8 py-8" ref={containerRef}>
            <div className="max-w-7xl mx-auto">
                {/* Header */}
                <div className="mb-8">
                    <div className="flex items-center gap-3 mb-2">
                        <History className="w-8 h-8 text-cyan-400" />
                        <h1 className="text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400">
                            My History
                        </h1>
                    </div>
                    <p className="text-slate-400">
                        {totalResults === 0
                            ? "No articles in your history yet"
                            : `Viewing ${articles.length} of ${totalResults} articles`}
                    </p>
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
                        onClose={handleCloseModal}
                        isAdmin={isAdmin}
                    />
                )}

                {/* Empty State */}
                {articles.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-20 bg-slate-900/30 rounded-2xl border border-blue-500/10">
                        <History className="w-16 h-16 text-slate-600 mb-4" />
                        <h3 className="text-xl font-semibold text-slate-400 mb-2">
                            No history yet
                        </h3>
                        <p className="text-slate-500 text-center max-w-md">
                            Articles you view and rate will appear here
                        </p>
                    </div>
                ) : (
                    <>
                        {/* Articles Grid */}
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
                            {articles.map((article) => (
                                <Article
                                    key={article.id}
                                    article={article}
                                    onSelect={() => handleSelectArticle(article)}
                                    isAdmin={isAdmin}
                                />
                            ))}
                        </div>

                        {/* Load More Button */}
                        {hasMore && (
                            <div className="flex justify-center">
                                <button
                                    onClick={handleLoadMore}
                                    disabled={loadingMore}
                                    className="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-500 to-cyan-500 text-white font-semibold rounded-lg hover:from-blue-600 hover:to-cyan-600 transition-all shadow-lg shadow-blue-500/25 disabled:opacity-50"
                                >
                                    {loadingMore ? (
                                        <>
                                            <RefreshCw className="w-5 h-5 animate-spin" />
                                            Loading...
                                        </>
                                    ) : (
                                        <>
                                            <ChevronDown className="w-5 h-5" />
                                            See More
                                        </>
                                    )}
                                </button>
                            </div>
                        )}

                        {/* End of Results */}
                        {!hasMore && articles.length > 0 && (
                            <div className="flex justify-center mt-8">
                                <p className="text-slate-500 text-sm">
                                    You've reached the end of your history
                                </p>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
};

export default HistoryPage;