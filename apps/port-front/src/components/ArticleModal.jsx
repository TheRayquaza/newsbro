import { useEffect, useState } from "react";
import api from "../api/api";
import { X, Calendar, ExternalLink, Trash2 } from "lucide-react";

const ArticleModal = ({ article, onClose, isAdmin }) => {
    const [feedback, setFeedback] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadFeedback = async () => {
            if (!article?.id) return;
            try {
                const data = await api.getArticleFeedback(article.id);
                console.log(data)
                setFeedback(data);
            } catch (error) {
                console.error("Failed to load article feedback", error);
            } finally {
                setLoading(false);
            }
        };
        if (isAdmin) loadFeedback();
        else setLoading(false)
    }, [article?.id, isAdmin]);

    const handleDelete = async () => {
        try {
            await api.deleteArticle(article.id);
            onClose();
        } catch (error) {
            alert("Failed to delete article");
            console.error("Failed to delete article", error);
        }
    };

    if (loading) return <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center"><div className="text-slate-300">Loading...</div></div>;

    return (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-slate-900/90 backdrop-blur-xl rounded-2xl border border-blue-500/20 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
                <div className="p-6">
                    <div className="flex items-start justify-between mb-4">
                        <h3 className="text-xl font-bold text-blue-400">{article.title}</h3>
                        <button onClick={onClose} className="text-slate-400 hover:text-slate-300">
                            <X className="w-6 h-6" />
                        </button>
                    </div>

                    <div className="flex flex-wrap items-center gap-2 text-xs mb-4">
                        <span className="px-3 py-1 bg-blue-500/20 text-blue-400 rounded-full">{article.category}</span>
                        {article.subcategory && (
                            <span className="px-3 py-1 bg-cyan-500/20 text-cyan-400 rounded-full">
                                {article.subcategory}
                            </span>
                        )}
                        {article.published_at && (
                            <span className="flex items-center gap-1 text-slate-500">
                                <Calendar className="w-3 h-3" />
                                {new Date(article.published_at).toLocaleDateString()}
                            </span>
                        )}
                    </div>

                    <p className="text-slate-300 mb-4">{article.abstract}</p>

                    {article.link && (
                        <a
                            href={article.link}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-2 text-blue-400 hover:text-cyan-400 transition mb-4"
                        >
                            <ExternalLink className="w-4 h-4" />
                            View Original Article
                        </a>
                    )}

                    {isAdmin && feedback && (
                        <div className="bg-slate-800/50 rounded-lg p-4 mb-4 border border-slate-700">
                            <h4 className="text-sm font-semibold text-slate-300 mb-3">Feedback Stats</h4>
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <p className="text-slate-400">Likes</p>
                                    <p className="text-lg font-bold text-green-400">{feedback.like_count}</p>
                                </div>
                                <div>
                                    <p className="text-slate-400">Dislikes</p>
                                    <p className="text-lg font-bold text-red-400">{feedback.dislike_count}</p>
                                </div>
                                <div>
                                    <p className="text-slate-400">Total</p>
                                    <p className="text-lg font-bold text-blue-400">{feedback.total_count}</p>
                                </div>
                                <div>
                                    <p className="text-slate-400">Like Ratio</p>
                                    <p className="text-lg font-bold text-cyan-400">{(feedback.like_ratio).toFixed(1)}%</p>
                                </div>
                            </div>
                        </div>
                    )}

                    {isAdmin && (
                        <div className="flex justify-end">
                            <button
                                onClick={handleDelete}
                                className="p-2 text-red-400 hover:bg-red-500/10 rounded-lg transition"
                            >
                                <Trash2 className="w-4 h-4" />
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default ArticleModal;
