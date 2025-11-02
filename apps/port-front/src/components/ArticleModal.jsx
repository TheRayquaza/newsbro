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

    if (loading) {
        return (
            <div 
                style={{ 
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    backdropFilter: 'blur(8px)'
                }} 
                className="fixed inset-0 flex items-center justify-center"
            >
                <div style={{ color: 'var(--nav-text)' }}>Loading...</div>
            </div>
        );
    }

    return (
        <div 
            style={{ 
                backgroundColor: 'rgba(0, 0, 0, 0.7)',
                backdropFilter: 'blur(8px)'
            }} 
            className="fixed inset-0 flex items-center justify-center p-4 z-50"
        >
            <div 
                style={{
                    backgroundColor: 'var(--nav-bg)',
                    borderColor: 'var(--nav-border)',
                    borderWidth: '1px',
                    borderStyle: 'solid'
                }}
                className="backdrop-blur-xl rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto"
            >
                <div className="p-6">
                    <div className="flex items-start justify-between mb-4">
                        <h3 
                            style={{ color: 'var(--nav-active-text)' }}
                            className="text-xl font-bold"
                        >
                            {article.title}
                        </h3>
                        <button 
                            onClick={onClose}
                            style={{ color: 'var(--nav-text-muted)' }}
                            onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-text)'}
                            onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-text-muted)'}
                            className="transition"
                        >
                            <X className="w-6 h-6" />
                        </button>
                    </div>

                    <div className="flex flex-wrap items-center gap-2 text-xs mb-4">
                        <span 
                            style={{
                                backgroundColor: 'var(--nav-active-bg)',
                                color: 'var(--nav-active-text)'
                            }}
                            className="px-3 py-1 rounded-full"
                        >
                            {article.category}
                        </span>
                        {article.subcategory && (
                            <span 
                                style={{
                                    backgroundColor: 'var(--search-button-active-bg)',
                                    color: 'var(--nav-active-text)'
                                }}
                                className="px-3 py-1 rounded-full"
                            >
                                {article.subcategory}
                            </span>
                        )}
                        {article.published_at && (
                            <span 
                                style={{ color: 'var(--nav-text-muted)' }}
                                className="flex items-center gap-1"
                            >
                                <Calendar className="w-3 h-3" />
                                {new Date(article.published_at).toLocaleDateString()}
                            </span>
                        )}
                    </div>

                    <p 
                        style={{ color: 'var(--nav-text)' }}
                        className="mb-4"
                    >
                        {article.abstract}
                    </p>

                    {article.link && (
                        <a
                            href={article.link}
                            target="_blank"
                            rel="noopener noreferrer"
                            style={{ color: 'var(--nav-active-text)' }}
                            onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-gradient-to)'}
                            onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
                            className="flex items-center gap-2 transition mb-4"
                        >
                            <ExternalLink className="w-4 h-4" />
                            View Original Article
                        </a>
                    )}

                    {isAdmin && feedback && (
                        <div 
                            style={{
                                backgroundColor: 'var(--nav-hover-bg)',
                                borderColor: 'var(--nav-border)',
                                borderWidth: '1px',
                                borderStyle: 'solid'
                            }}
                            className="rounded-lg p-4 mb-4"
                        >
                            <h4 
                                style={{ color: 'var(--nav-text)' }}
                                className="text-sm font-semibold mb-3"
                            >
                                Feedback Stats
                            </h4>
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <p style={{ color: 'var(--nav-text-muted)' }}>Likes</p>
                                    <p 
                                        style={{ color: 'var(--success)' }}
                                        className="text-lg font-bold"
                                    >
                                        {feedback.like_count}
                                    </p>
                                </div>
                                <div>
                                    <p style={{ color: 'var(--nav-text-muted)' }}>Dislikes</p>
                                    <p 
                                        style={{ color: 'var(--error)' }}
                                        className="text-lg font-bold"
                                    >
                                        {feedback.dislike_count}
                                    </p>
                                </div>
                                <div>
                                    <p style={{ color: 'var(--nav-text-muted)' }}>Total</p>
                                    <p 
                                        style={{ color: 'var(--nav-active-text)' }}
                                        className="text-lg font-bold"
                                    >
                                        {feedback.total_count}
                                    </p>
                                </div>
                                <div>
                                    <p style={{ color: 'var(--nav-text-muted)' }}>Like Ratio</p>
                                    <p 
                                        style={{ color: 'var(--nav-gradient-to)' }}
                                        className="text-lg font-bold"
                                    >
                                        {(feedback.like_ratio).toFixed(1)}%
                                    </p>
                                </div>
                            </div>
                        </div>
                    )}

                    {isAdmin && (
                        <div className="flex justify-end">
                            <button
                                onClick={handleDelete}
                                style={{ 
                                    color: 'var(--error)',
                                    backgroundColor: 'transparent'
                                }}
                                onMouseEnter={(e) => {
                                    e.currentTarget.style.backgroundColor = 'var(--error-hover)';
                                    e.currentTarget.style.color = 'var(--error)';
                                }}
                                onMouseLeave={(e) => {
                                    e.currentTarget.style.backgroundColor = 'transparent';
                                    e.currentTarget.style.color = 'var(--error)';
                                }}
                                className="p-2 rounded-lg transition"
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
