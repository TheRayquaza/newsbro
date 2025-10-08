import { X, ThumbsUp, ThumbsDown, ExternalLink, Calendar } from "lucide-react";

export default function ArticleModal({ article, onClose, feedback, onFeedback }) {
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
                    {(
                        <div className="flex items-center gap-3 mt-4">
                            <button
                                onClick={() => onFeedback(true)}
                                className={`flex items-center gap-1 px-3 py-2 rounded-lg transition-all ${feedback === true
                                    ? "text-green-400 hover:text-green-300 hover:bg-green-500/10"
                                    : "bg-green-500 text-white border border-green-600 shadow-lg scale-105"
                                    }`}
                            >
                                <ThumbsUp className="w-4 h-4" />
                                {feedback === true ? "Liked" : "Like"}
                            </button>
                            <button
                                onClick={() => onFeedback(false)}
                                className={`flex items-center gap-1 px-3 py-2 rounded-lg transition-all ${feedback === false
                                    ? "text-red-400 hover:text-red-300 hover:bg-red-500/10"
                                    : "bg-red-500 text-white border border-red-600 shadow-lg scale-105"
                                    }`}
                            >
                                <ThumbsDown className="w-4 h-4" />
                                {feedback === false ? "Disliked" : "Dislike"}
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
