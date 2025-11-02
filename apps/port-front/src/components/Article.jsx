import React, { useContext, useEffect, useState } from "react";
import { Calendar, ThumbsUp, ThumbsDown, ExternalLink, Link } from "lucide-react";
import api from "../api/api";
import { AuthContext } from "../contexts/Auth";

const Article = ({ article, onSelect, isAdmin }) => {
    const [userFeedback, setUserFeedback] = useState(null);
    const { token } = useContext(AuthContext);

    useEffect(() => {
        if (article?.liked_value !== -1)
            setUserFeedback(article.liked_value === 1 ? true : false);
    }, [article]);

    const handleFeedback = async (value) => {
        if (!article?.id) return;

        try {
            if (userFeedback === value) {
                await api.deleteArticleFeedback(article.id, token);
                setUserFeedback(null);
            } else {
                await api.createArticleFeedback(article.id, { value }, token);
                setUserFeedback(value);
            }
        } catch (error) {
            alert("Failed to submit feedback");
            console.error("Failed to submit feedback", error);
        }
    };

    return (
        <div
            style={{
                backgroundColor: 'var(--nav-bg)',
                borderColor: 'var(--nav-border)',
                borderWidth: '1px',
                borderStyle: 'solid'
            }}
            className="backdrop-blur-xl rounded-xl transition-all p-6 cursor-pointer"
            onClick={onSelect}
            onMouseEnter={(e) => e.currentTarget.style.borderColor = 'var(--nav-active-text)'}
            onMouseLeave={(e) => e.currentTarget.style.borderColor = 'var(--nav-border)'}
        >
            <h3 
                style={{ color: 'var(--nav-text)' }}
                className="text-lg font-semibold mb-2 line-clamp-2"
            >
                {article.title}
            </h3>
            <p 
                style={{ color: 'var(--nav-text-muted)' }}
                className="text-sm mb-3 line-clamp-3"
            >
                {article.abstract}
            </p>

            <div className="flex flex-wrap items-center gap-2 text-xs mb-3">
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

            {article.link && (
                <a
                    href={article.link}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{ color: 'var(--nav-active-text)' }}
                    onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-gradient-to)'}
                    onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
                    className="flex items-center gap-1 text-sm mb-3 transition"
                    onClick={(e) => e.stopPropagation()}
                >
                    <ExternalLink className="w-4 h-4" />
                    View Article
                </a>
            )}

            {article.rss_link && (
                <a
                    href={article.rss_link}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{ color: 'var(--nav-active-text)' }}
                    onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-gradient-to)'}
                    onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
                    className="flex items-center gap-1 text-sm mb-3 transition"
                    onClick={(e) => e.stopPropagation()}
                >
                    <Link className="w-4 h-4" />
                    RSS
                </a>
            )}

            <div className="flex items-center gap-3">
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        handleFeedback(true);
                    }}
                    style={{
                        backgroundColor: userFeedback === true ? 'var(--success)' : 'transparent',
                        color: userFeedback === true ? '#ffffff' : 'var(--success)',
                        borderColor: 'var(--success)',
                        borderWidth: '1px',
                        borderStyle: 'solid'
                    }}
                    onMouseEnter={(e) => {
                        if (userFeedback === true) {
                            e.currentTarget.style.backgroundColor = 'var(--success-hover)';
                        } else {
                            e.currentTarget.style.backgroundColor = 'var(--success-bg-hover)';
                        }
                    }}
                    onMouseLeave={(e) => {
                        if (userFeedback === true) {
                            e.currentTarget.style.backgroundColor = 'var(--success)';
                        } else {
                            e.currentTarget.style.backgroundColor = 'transparent';
                        }
                    }}
                    className="flex items-center gap-1 px-2 py-1 rounded-lg transition-all"
                >
                    <ThumbsUp className="w-4 h-4" />
                    {isAdmin && article.likes}
                </button>

                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        handleFeedback(false);
                    }}
                    style={{
                        backgroundColor: userFeedback === false ? 'var(--error)' : 'transparent',
                        color: userFeedback === false ? '#ffffff' : 'var(--error)',
                        borderColor: 'var(--error)',
                        borderWidth: '1px',
                        borderStyle: 'solid'
                    }}
                    onMouseEnter={(e) => {
                        if (userFeedback === false) {
                            e.currentTarget.style.backgroundColor = 'var(--error-hover)';
                        } else {
                            e.currentTarget.style.backgroundColor = 'var(--error-bg-hover)';
                        }
                    }}
                    onMouseLeave={(e) => {
                        if (userFeedback === false) {
                            e.currentTarget.style.backgroundColor = 'var(--error)';
                        } else {
                            e.currentTarget.style.backgroundColor = 'transparent';
                        }
                    }}
                    className="flex items-center gap-1 px-2 py-1 rounded-lg transition-all"
                >
                    <ThumbsDown className="w-4 h-4" />
                    {isAdmin && article.dislikes}
                </button>
            </div>
        </div>
    );
};

export default Article;