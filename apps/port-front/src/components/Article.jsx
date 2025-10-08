import React, { useContext, useEffect, useState } from "react";
import { Calendar, Trash2, ThumbsUp, ThumbsDown, ExternalLink } from "lucide-react";
import api from "../api/api";
import { AuthContext } from "../contexts/Auth";

const Article = ({ article, onSelect, isAdmin, onDelete }) => {
    const [likes, setLikes] = useState(article.likes || 0);
    const [dislikes, setDislikes] = useState(article.dislikes || 0);
    const [userFeedback, setUserFeedback] = useState(null); // true = liked, false = disliked, null = none
    const [loading, setLoading] = useState(true);
    const { token } = useContext(AuthContext);

    // Load stats if admin (optional)
    useEffect(() => {
        const loadStats = async () => {
            if (!article?.id) return;
            try {
                // Placeholder until API implemented
                const stats = { likes: article.likes || 0, dislikes: article.dislikes || 0, user_feedback: null };
                setLikes(stats.likes);
                setDislikes(stats.dislikes);
                setUserFeedback(stats.user_feedback);
            } catch (error) {
                console.error("Failed to load article stats", error);
            } finally {
                setLoading(false);
            }
        };

        if (isAdmin) loadStats();
        else setLoading(false);
    }, [article, isAdmin]);

    const handleFeedback = async (value) => {
        if (!article?.id) return;

        try {
            if (userFeedback === value) {
                // Remove feedback
                await api.deleteArticleFeedback(article.id, token);
                setUserFeedback(null);
            } else {
                // Submit new feedback
                await api.createArticleFeedback(article.id, { value }, token);
                setUserFeedback(value);
            }
        } catch (error) {
            alert("Failed to submit feedback");
            console.error("Failed to submit feedback", error);
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div
            className="bg-slate-900/60 backdrop-blur-xl rounded-xl border border-blue-500/20 hover:border-blue-500/40 transition-all p-6 cursor-pointer"
            onClick={onSelect}
        >
            <h3 className="text-lg font-semibold text-slate-200 mb-2 line-clamp-2">{article.title}</h3>
            <p className="text-slate-400 text-sm mb-3 line-clamp-3">{article.abstract}</p>

            <div className="flex flex-wrap items-center gap-2 text-xs mb-3">
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

            {/* Link */}
            {article.link && (
                <a
                    href={article.link}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-1 text-blue-400 text-sm hover:text-cyan-400 mb-3"
                    onClick={(e) => e.stopPropagation()}
                >
                    <ExternalLink className="w-4 h-4" />
                    View Article
                </a>
            )}

            <div className="flex items-center gap-3">
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        handleFeedback(true);
                    }}
                    className={`flex items-center gap-1 px-2 py-1 rounded-lg transition-all ${userFeedback === true
                            ? "text-green-400 hover:text-green-300 hover:bg-green-500/10"
                            : "bg-green-500 text-white border border-green-600 shadow-lg scale-105"
                        }`}
                >
                    <ThumbsUp className="w-4 h-4" />
                    {isAdmin ? likes : null}
                </button>
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        handleFeedback(false);
                    }}
                    className={`flex items-center gap-1 px-2 py-1 rounded-lg transition-all ${userFeedback === false
                            ? "text-red-400 hover:text-red-300 hover:bg-red-500/10"
                            : "bg-red-500 text-white border border-red-600 shadow-lg scale-105"
                        }`}
                >
                    <ThumbsDown className="w-4 h-4" />
                    {isAdmin ? dislikes : null}
                </button>
            </div>


            {/* Admin actions */}
            {isAdmin && (
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        onDelete(article);
                    }}
                    className="p-2 text-red-400 hover:bg-red-500/10 rounded-lg transition"
                >
                    <Trash2 className="w-4 h-4" />
                </button>
            )}
        </div>
    );
};

export default Article;
