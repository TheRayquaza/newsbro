import React, { useContext, useEffect, useState, useRef } from 'react';
import { Link, ChevronDown } from 'lucide-react';
import { AuthContext } from '../contexts/Auth';
import DashboardBreadcrumb from '../components/DashboardBreadcrumb';
import Article from '../components/Article';
import api from '../api/api';
import RSS from '../components/RSS';
import ArticleModal from '../components/ArticleModal';

const Home = () => {
  const { isAdmin } = useContext(AuthContext);

  // UI State
  const [breadcrumb, setBreadcrumb] = useState([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  // Data State
  const [articleFeedback, setArticleFeedback] = useState(null);
  const [selectedArticle, setSelectedArticle] = useState(null);
  const [currentItems, setCurrentItems] = useState([]);
  const [articles, setArticles] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);

  // Caching State
  const [articlesCache, setArticlesCache] = useState({});
  const [feedbackCache, setFeedbackCache] = useState({});
  const rssTreeCache = useRef(null);

  const pendingRequests = useRef({
    counts: {},
    articles: {},
    feedback: {}
  });

  const ARTICLES_PER_PAGE = 20;

  useEffect(() => {
    if (breadcrumb.length === 0) {
      if (rssTreeCache.current) {
        setCurrentItems(rssTreeCache.current);
        return;
      }

      setLoading(true);
      api.getRssTree().then(response => {
        rssTreeCache.current = response;
        setCurrentItems(response);
        setLoading(false);
      }).catch(() => {
        setCurrentItems([]);
        setLoading(false);
      });
    }
  }, [breadcrumb]);

  const getStartOfWeek = () => {
    const now = new Date();
    const dayOfWeek = now.getDay();
    const diff = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
    const startOfWeek = new Date(now);
    startOfWeek.setDate(now.getDate() - diff);
    startOfWeek.setHours(0, 0, 0, 0);
    return startOfWeek;
  };

  const getEndOfWeek = () => {
    const startOfWeek = getStartOfWeek();
    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6);
    endOfWeek.setHours(23, 59, 59, 999);
    return endOfWeek;
  };

  const loadArticlesInFeed = async (feedName, page = 1, append = false) => {
    const cacheKey = `${feedName}_page_${page}`;
    
    if (articlesCache[cacheKey] && !append) {
      setArticles(articlesCache[cacheKey]);
      setHasMore(articlesCache[cacheKey].length === ARTICLES_PER_PAGE);
      return { articles: articlesCache[cacheKey] };
    }
  
    if (pendingRequests.current.articles[cacheKey]) {
      return pendingRequests.current.articles[cacheKey];
    }
  
    if (append) {
      setLoadingMore(true);
    } else {
      setLoading(true);
    }
  
    const startOfWeek = getStartOfWeek();
    const endOfWeek = getEndOfWeek();
  
    const requestPromise = api.getArticles({
      "feed_name": feedName,
      "start_date": startOfWeek.toISOString(),
      "end_date": endOfWeek.toISOString(),
      "limit": ARTICLES_PER_PAGE,
      "offset": (page - 1) * ARTICLES_PER_PAGE
    })
      .then(response => {
        const newArticles = response.articles || [];
        
        const sortedArticles = newArticles.sort((a, b) => {
          const dateA = new Date(a.published_at);
          const dateB = new Date(b.published_at);
          return dateB - dateA;
        });
        
        if (append) {
          setArticles(prev => [...prev, ...sortedArticles]);
        } else {
          setArticles(sortedArticles);
        }
        
        setArticlesCache(prev => ({ ...prev, [cacheKey]: sortedArticles }));
        setHasMore(sortedArticles.length === ARTICLES_PER_PAGE);
        setLoadingMore(false);
        setLoading(false);
        delete pendingRequests.current.articles[cacheKey];
        return response;
      })
      .catch(error => {
        alert("Failed to load articles for the selected feed: " + error.message);
        setLoadingMore(false);
        setLoading(false);
        delete pendingRequests.current.articles[cacheKey];
        return { articles: [] };
      });
  
    pendingRequests.current.articles[cacheKey] = requestPromise;
    return requestPromise;
  };

  const handleLoadMore = async () => {
    if (loadingMore || !hasMore) return;
    
    const currentFeedName = breadcrumb[breadcrumb.length - 1]?.name;
    if (!currentFeedName) return;

    const nextPage = currentPage + 1;
    setCurrentPage(nextPage);
    await loadArticlesInFeed(currentFeedName, nextPage, true);
  };

  const handleItemClick = async (item) => {
    setLoading(true);

    if (item.children && item.children.length > 0) {
      setCurrentItems(item.children);
      setBreadcrumb([...breadcrumb, item]);
      setArticles([]);
      setCurrentPage(1);
      setHasMore(true);
      setLoading(false);
    }
    else if (item.link) {
      setCurrentPage(1);
      setHasMore(true);
      await loadArticlesInFeed(item.name, 1, false);
      setBreadcrumb([...breadcrumb, item]);
      setCurrentItems([]);
      setLoading(false);
    } else {
      setLoading(false);
    }
  };

  const handleViewArticle = async (article) => {
    setSelectedArticle(article);
    if (isAdmin) {
      if (feedbackCache[article.id] !== undefined) {
        setArticleFeedback(feedbackCache[article.id]);
        return;
      }

      if (pendingRequests.current.feedback[article.id]) {
        try {
          const feedback = await pendingRequests.current.feedback[article.id];
          setArticleFeedback(feedback);
        } catch (err) {
          console.error("Failed to load feedback:", err);
        }
        return;
      }

      const feedbackPromise = api.getArticleFeedback(article.id)
        .then(feedback => {
          setFeedbackCache(prev => ({ ...prev, [article.id]: feedback }));
          setArticleFeedback(feedback);
          delete pendingRequests.current.feedback[article.id];
          return feedback;
        })
        .catch(err => {
          console.error("Failed to load feedback:", err);
          delete pendingRequests.current.feedback[article.id];
          throw err;
        });

      pendingRequests.current.feedback[article.id] = feedbackPromise;
    }
  };

  const handleDeleteArticle = async (id) => {
    if (!confirm("Are you sure you want to delete this article?")) return;
    try {
      await api.deleteArticle(id);
      const updatedArticles = articles.filter((a) => a.id !== id);
      setArticles(updatedArticles);

      const currentFeedName = breadcrumb[breadcrumb.length - 1]?.name;
      if (currentFeedName) {
        const cacheKey = `${currentFeedName}_page_${currentPage}`;
        setArticlesCache(prev => ({ ...prev, [cacheKey]: updatedArticles }));
      }

      setSelectedArticle(null);
    } catch {
      alert("Failed to delete article.");
    }
  };

  const getFeedCount = (item) => {
    if (!item.children || item.children.length === 0) return 0;

    let feedCount = 0;
    const countFeeds = (items) => {
      items.forEach(child => {
        if (child.link && child.link !== "") {
          feedCount++;
        }
        if (child.children && child.children.length > 0) {
          countFeeds(child.children);
        }
      });
    };

    countFeeds(item.children);
    return feedCount;
  };

  const handleFeedbackSubmit = async (isPositive) => {
    try {
      await api.createArticleFeedback(selectedArticle.id, { value: isPositive });
      setArticleFeedback(isPositive);
      setFeedbackCache(prev => ({ ...prev, [selectedArticle.id]: isPositive }));
    } catch (err) {
      console.error("Failed to submit feedback:", err);
    }
  };

  const isArticlesView = articles.length > 0;
  const currentLevel = breadcrumb.length > 0 ? breadcrumb[breadcrumb.length - 1] : null;

  return (
    <div>
      {/* Header */}
      <div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex flex-col w-full gap-2">
            <div className="flex items-center justify-between">
              <h1 
                style={{
                  background: 'linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                  backgroundClip: 'text'
                }}
                className="text-3xl font-bold"
              >
                News Feed Hub
              </h1>
            </div>
            <h2 style={{ color: 'var(--nav-text-muted)' }}>
              Curated news articles from various RSS feeds, updated weekly.
            </h2>
            <DashboardBreadcrumb
              breadcrumb={breadcrumb}
              setBreadcrumb={setBreadcrumb}
              setCurrentItems={setCurrentItems}
              setArticles={setArticles}
              loadArticlesInFeed={(feedName) => {
                setCurrentPage(1);
                setHasMore(true);
                return loadArticlesInFeed(feedName, 1, false);
              }}
            />
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Feeds Grid View */}
        {!isArticlesView && currentItems.length > 0 && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {currentItems.map((item) => (
                <RSS
                  key={item.name}
                  item={item}
                  handleItemClick={handleItemClick}
                  getFeedCount={getFeedCount}
                />
              ))}
            </div>
          </div>
        )}

        {/* Article Modal */}
        {selectedArticle && (
          <ArticleModal
            article={selectedArticle}
            onClose={() => setSelectedArticle(null)}
            isAdmin={isAdmin}
            feedback={articleFeedback}
            onFeedback={handleFeedbackSubmit}
          />
        )}

        {/* Articles View */}
        {isArticlesView && (
          <div className="space-y-6">
            <div 
              style={{
                background: 'var(--nav-active-bg)',
                borderColor: 'var(--nav-border)',
                borderWidth: '1px',
                borderStyle: 'solid'
              }}
              className="rounded-xl p-6"
            >
              <h2 
                style={{ color: 'var(--nav-active-text)' }}
                className="text-2xl font-bold mb-2"
              >
                {currentLevel?.display_name}
              </h2>
              <p style={{ color: 'var(--nav-text-muted)' }} className="text-sm">
                {currentLevel?.description}
              </p>
              {currentLevel?.link && (
                <a
                  href={currentLevel.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  style={{ color: 'var(--success)' }}
                  onMouseEnter={(e) => e.currentTarget.style.color = 'var(--success-hover)'}
                  onMouseLeave={(e) => e.currentTarget.style.color = 'var(--success)'}
                  className="inline-block mt-2 text-xs hover:underline transition"
                >
                  <Link className="inline-block align-middle mr-2" size={14} /> RSS Feed
                </a>
              )}
            </div>

            {/* Articles List */}
            <div className="space-y-4">
              {loading ? (
                <div className="flex items-center justify-center py-12">
                  <div 
                    style={{ borderTopColor: 'var(--nav-active-text)' }}
                    className="animate-spin rounded-full h-8 w-8 border-t-2"
                  ></div>
                </div>
              ) : articles.length > 0 ? (
                <>
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
                  
                  {/* Load More Button */}
                  {hasMore && (
                    <div className="flex justify-center pt-6">
                      <button
                        onClick={handleLoadMore}
                        disabled={loadingMore}
                        style={{
                          backgroundColor: loadingMore ? 'var(--nav-hover-bg)' : 'var(--nav-active-bg)',
                          color: 'var(--nav-active-text)',
                          borderColor: 'var(--nav-border)',
                          borderWidth: '1px',
                          borderStyle: 'solid'
                        }}
                        onMouseEnter={(e) => {
                          if (!loadingMore) {
                            e.currentTarget.style.backgroundColor = 'var(--nav-active-text)';
                            e.currentTarget.style.color = '#ffffff';
                          }
                        }}
                        onMouseLeave={(e) => {
                          if (!loadingMore) {
                            e.currentTarget.style.backgroundColor = 'var(--nav-active-bg)';
                            e.currentTarget.style.color = 'var(--nav-active-text)';
                          }
                        }}
                        className="flex items-center gap-2 px-6 py-3 rounded-lg font-medium transition-all"
                      >
                        {loadingMore ? (
                          <>
                            <div 
                              style={{ borderTopColor: 'var(--nav-active-text)' }}
                              className="animate-spin rounded-full h-4 w-4 border-t-2"
                            ></div>
                            Loading...
                          </>
                        ) : (
                          <>
                            Load More
                            <ChevronDown className="w-4 h-4" />
                          </>
                        )}
                      </button>
                    </div>
                  )}
                </>
              ) : (
                <div className="text-center py-12">
                  <p style={{ color: 'var(--nav-text-muted)' }}>
                    No articles found for this week.
                  </p>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Empty State */}
        {!isArticlesView && currentItems.length === 0 && (
          <div className="text-center py-12">
            <p style={{ color: 'var(--nav-text-muted)' }}>No items available.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
