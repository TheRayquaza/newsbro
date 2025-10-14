import React, { useContext, useEffect, useState } from 'react';
import { ChevronLeft, Search, Bookmark, Share2, ExternalLink, Home } from 'lucide-react';
import { AuthContext } from '../contexts/Auth';
import api from '../api/api';

const Dashboard = () => {
  const { isAdmin } = useContext(AuthContext);
  const [breadcrumb, setBreadcrumb] = useState([]);
  const [currentItems, setCurrentItems] = useState([]);
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedArticle, setSelectedArticle] = useState(null);

  // Load root categories on mount
  useEffect(() => {
    if (breadcrumb.length === 0) {
      setLoading(true);
      api.getRssTree().then(response => {
        setCurrentItems(response);
        setLoading(false);
      }).catch(() => {
        setCurrentItems([]);
        setLoading(false);
      });
    }
  }, [breadcrumb]);

  const handleItemClick = async (item) => {
    setLoading(true);
    setSearchTerm('');

    if (item.children && item.children.length > 0) {
      setCurrentItems(item.children);
      setBreadcrumb([...breadcrumb, item]);
      setArticles([]);
    }
    else if (item.link) {
      setTimeout(() => {
        const mockArticles = generateArticlesForFeed(item.name);
        setArticles(mockArticles);
        setBreadcrumb([...breadcrumb, item]);
        setCurrentItems([]);
        setLoading(false);
      }, 300);
      return;
    }

    setLoading(false);
  };

  const handleBack = () => {
    if (breadcrumb.length === 0) return;

    const newBreadcrumb = breadcrumb.slice(0, -1);
    setBreadcrumb(newBreadcrumb);
    setArticles([]);
    setSearchTerm('');

    if (newBreadcrumb.length !== 0) {
      const parent = newBreadcrumb[newBreadcrumb.length - 1];
      setCurrentItems(parent.children || []);
    }
  };

  const filteredArticles = articles.filter(article =>
    article.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    article.excerpt.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const isArticlesView = articles.length > 0;
  const currentLevel = breadcrumb.length > 0 ? breadcrumb[breadcrumb.length - 1] : null;

  return (
    <div>
      {/* Header */}
      <div className="sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex flex-col w-full gap-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                {breadcrumb.length > 0 && (
                  <button
                    onClick={handleBack}
                    className="p-2 hover:bg-slate-800/60 rounded-lg transition"
                  >
                    <ChevronLeft className="text-blue-400" size={22} />
                  </button>
                )}
                <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
                  News Feed Hub
                </h1>
              </div>
            </div>

            {/* Breadcrumb under title */}
            <ol className="flex flex-wrap items-center gap-2 text-sm mt-1">
              <li key={""} className="inline-flex items-center">
                <a
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setBreadcrumb([]);
                  }}
                  className={`flex items-center px-3 py-1.5 rounded-lg transition-all duration-200 ${breadcrumb.length <= 1
                    ? 'text-blue-400 font-semibold bg-blue-500/10 cursor-default'
                    : 'text-slate-400 hover:text-blue-300 hover:bg-blue-500/10'
                    }`}
                >
                  <Home className="w-5 h-5 text-blue-400" />
                </a>
              </li>
              {
                breadcrumb.map((b, index) => {
                  const isLast = index === breadcrumb.length - 1;
                  return (
                    <li key={b.name} className="inline-flex items-center">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="w-4 h-4 mx-2 text-slate-500"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M9 5l7 7-7 7"
                        />
                      </svg>
                      <a
                        href="#"
                        onClick={(e) => {
                          e.preventDefault();
                          const newBreadcrumb = breadcrumb.slice(0, index + 1);
                          setBreadcrumb(newBreadcrumb);
                          if (b.children && b.children.length > 0) {
                            setCurrentItems(b.children);
                            setArticles([]);
                          } else {
                            const mockArticles = generateArticlesForFeed(b.name);
                            setArticles(mockArticles);
                            setCurrentItems([]);
                          }
                        }}
                        className={`flex items-center px-3 py-1.5 rounded-lg transition-all duration-200 ${isLast
                          ? 'text-blue-400 font-semibold bg-blue-500/10 cursor-default'
                          : 'text-slate-400 hover:text-blue-300 hover:bg-blue-500/10'
                          }`}
                      >
                        {b.name}
                      </a>
                    </li>
                  );
                }
                )}
            </ol>
          </div>
        </div>
      </div>


      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Feeds Grid View */}
        {!isArticlesView && currentItems.length > 0 && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {currentItems.map((item) => (
                <button
                  key={item.name}
                  onClick={() => handleItemClick(item)}
                  className="group relative overflow-hidden bg-gradient-to-br from-slate-800/40 to-slate-800/20 border border-blue-500/20 rounded-2xl p-8 hover:border-blue-500/50 transition-all duration-300 text-left"
                >
                  <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-cyan-500/10 opacity-0 group-hover:opacity-100 transition-opacity" />
                  <div className="relative z-10">
                    <h3 className="text-2xl font-bold text-blue-300 group-hover:text-blue-200 transition mb-2">
                      {item.display_name}
                    </h3>
                    <p className="text-slate-400 group-hover:text-slate-300 transition mb-6 text-sm line-clamp-2">
                      {item.description}
                    </p>
                    <div className="flex items-center justify-between">
                      <div className="inline-block px-4 py-2 bg-blue-500/20 group-hover:bg-blue-500/30 text-blue-300 rounded-lg text-sm font-medium transition">
                        {item.children?.length || 0} {item.children?.length === 1 ? 'item' : 'items'}
                      </div>

                      {item.link ? (
                        <span className="ml-2 px-3 py-1 text-xs rounded-full bg-green-500/20 text-green-300 border border-green-500/30">
                          RSS
                        </span>
                      ) : (
                        <span className="ml-2 px-3 py-1 text-xs rounded-full bg-purple-500/20 text-purple-300 border border-purple-500/30">
                          Group
                        </span>
                      )}
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </div>
        )}

        {/* Articles View */}
        {isArticlesView && (
          <div className="space-y-6">
            {/* Feed Header */}
            <div className="bg-gradient-to-r from-blue-500/10 to-cyan-500/10 border border-blue-500/20 rounded-xl p-6">
              <h2 className="text-2xl font-bold text-blue-300 mb-2">{currentLevel?.name}</h2>
              <p className="text-slate-400 text-sm">
                {currentLevel?.description}
              </p>
            </div>

            {/* Search Bar */}
            <div className="relative">
              <Search className="absolute left-3 top-3 text-slate-400" size={20} />
              <input
                type="text"
                placeholder="Search articles..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full pl-10 pr-4 py-3 bg-slate-800/40 border border-blue-500/20 rounded-xl text-slate-200 placeholder-slate-500 focus:outline-none focus:border-blue-500/50 transition"
              />
            </div>

            {/* Articles List */}
            <div className="space-y-4">
              {loading ? (
                <div className="flex items-center justify-center py-12">
                  <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-blue-400"></div>
                </div>
              ) : filteredArticles.length > 0 ? (
                filteredArticles.map(article => (
                  <article
                    key={article.id}
                    onClick={() => setSelectedArticle(article)}
                    className="group bg-slate-800/40 border border-blue-500/20 rounded-xl p-6 hover:bg-slate-800/60 hover:border-blue-500/40 transition duration-300 backdrop-blur-sm cursor-pointer"
                  >
                    <div className="space-y-4">
                      <div>
                        <h3 className="text-lg font-semibold text-blue-300 group-hover:text-blue-200 transition mb-2">
                          {article.title}
                        </h3>
                        <p className="text-slate-400 text-sm">
                          {article.excerpt}
                        </p>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-slate-500">
                          {new Date(article.date).toLocaleDateString('en-US', {
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric'
                          })}
                        </span>
                        <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition">
                          <button className="flex items-center gap-1 px-3 py-1.5 text-xs bg-blue-500/20 hover:bg-blue-500/30 text-blue-300 rounded transition">
                            <Bookmark size={14} />
                          </button>
                          <button className="flex items-center gap-1 px-3 py-1.5 text-xs bg-slate-700/40 hover:bg-slate-700/60 text-slate-300 rounded transition">
                            <Share2 size={14} />
                          </button>
                          <button className="flex items-center gap-1 px-3 py-1.5 text-xs bg-slate-700/40 hover:bg-slate-700/60 text-slate-300 rounded transition">
                            <ExternalLink size={14} />
                          </button>
                        </div>
                      </div>
                    </div>
                  </article>
                ))
              ) : (
                <div className="text-center py-12">
                  <p className="text-slate-400">No articles found matching your search.</p>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Empty State */}
        {!isArticlesView && currentItems.length === 0 && (
          <div className="text-center py-12">
            <p className="text-slate-400">No items available.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;