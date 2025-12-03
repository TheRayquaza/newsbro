import React, { useContext, useEffect, useState, useRef, useCallback } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Link, ChevronDown } from 'lucide-react';
import { AuthContext } from '../contexts/Auth';
import DashboardBreadcrumb from '../components/DashboardBreadcrumb';
import Article from '../components/Article';
import api from '../api/api';
import RSS from '../components/RSS';
import ArticleModal from '../components/ArticleModal';

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

const Home = () => {
  const { isAdmin } = useContext(AuthContext);
  const navigate = useNavigate();
  const { '*': urlPath } = useParams();

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
  const pendingRequests = useRef({ counts: {}, articles: {}, feedback: {} });

  const ARTICLES_PER_PAGE = 21;

  const findItemByPath = (tree, pathSegments) => {
    if (!tree || pathSegments.length === 0) return null;

    let current = tree;
    const breadcrumbPath = [];

    for (const segment of pathSegments) {
      const decodedSegment = decodeURIComponent(segment);
      const found = current.find(item => item.name === decodedSegment);

      if (!found) return null;

      breadcrumbPath.push(found);

      if (found.children && found.children.length > 0) {
        current = found.children;
      } else {
        return { item: found, breadcrumb: breadcrumbPath };
      }
    }

    return {
      item: breadcrumbPath[breadcrumbPath.length - 1],
      breadcrumb: breadcrumbPath
    };
  };

  const buildUrlPath = (breadcrumbItems) => {
    if (breadcrumbItems.length === 0) return '/';
    return '/' + breadcrumbItems.map(item => encodeURIComponent(item.name)).join('/');
  };

  const loadArticlesInFeed = useCallback(
    async (feedName, page = 1, append = false) => {
      const cacheKey = `${feedName}_page_${page}`;

      if (articlesCache[cacheKey] && !append) {
        setArticles(articlesCache[cacheKey]);
        setHasMore(articlesCache[cacheKey].length === ARTICLES_PER_PAGE);
        return { articles: articlesCache[cacheKey] };
      }

      if (pendingRequests.current.articles[cacheKey]) {
        return pendingRequests.current.articles[cacheKey];
      }

      if (append) setLoadingMore(true);
      else setLoading(true);

      const startOfWeek = getStartOfWeek();
      const endOfWeek = getEndOfWeek();

      const requestPromise = api.getArticles({
        feed_name: feedName,
        start_date: startOfWeek.toISOString(),
        end_date: endOfWeek.toISOString(),
        limit: ARTICLES_PER_PAGE,
        offset: (page - 1) * ARTICLES_PER_PAGE,
      })
        .then(response => {
          const newArticles = response.articles || [];
          const sortedArticles = newArticles.sort(
            (a, b) => new Date(b.published_at) - new Date(a.published_at)
          );

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
    },
    [
      articlesCache,
      setArticles,
      setArticlesCache,
      setHasMore,
      setLoading,
      setLoadingMore,
      pendingRequests
    ]
  );

  useEffect(() => {
    const loadInitialData = async () => {
      setLoading(true);

      try {
        if (!rssTreeCache.current) {
          const response = await api.getRssTree();
          rssTreeCache.current = response;
        }

        if (urlPath) {
          const pathSegments = urlPath.split('/').filter(Boolean);
          const result = findItemByPath(rssTreeCache.current, pathSegments);

          if (result) {
            const { item, breadcrumb: foundBreadcrumb } = result;
            setBreadcrumb(foundBreadcrumb);

            if (item.link) {
              await loadArticlesInFeed(item.name, 1, false);
              setCurrentItems([]);
            } else if (item.children) {
              setCurrentItems(item.children);
              setArticles([]);
            }
          } else {
            navigate('/', { replace: true });
            setCurrentItems(rssTreeCache.current);
            setBreadcrumb([]);
          }
        } else {
          setCurrentItems(rssTreeCache.current);
          setBreadcrumb([]);
          setArticles([]);
        }
      } catch (error) {
        console.error('Failed to load RSS tree:', error);
        setCurrentItems([]);
      } finally {
        setLoading(false);
      }
    };

    loadInitialData();
  }, [urlPath, navigate, loadArticlesInFeed]);

  const handleLoadMore = async () => {
    if (loadingMore || !hasMore) return;

    const currentFeedName = breadcrumb[breadcrumb.length - 1]?.name;
    if (!currentFeedName) return;

    const nextPage = currentPage + 1;
    setCurrentPage(nextPage);
    await loadArticlesInFeed(currentFeedName, nextPage, true);
  };

  const handleItemClick = async (item) => {
    const newBreadcrumb = [...breadcrumb, item];
    const newPath = buildUrlPath(newBreadcrumb);

    navigate(newPath);
  };

  const handleBreadcrumbClick = (index) => {
    if (index === -1) {
      navigate('/');
    } else {
      const newBreadcrumb = breadcrumb.slice(0, index + 1);
      const newPath = buildUrlPath(newBreadcrumb);
      navigate(newPath);
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

  return (
    <div style={{ padding: '2rem', marginLeft: '17%', marginRight: '17%' }}>
      {/* Header */}
      <div style={{ marginBottom: '2rem' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
          News Feed Hub
        </h1>
        <p style={{ color: 'var(--text-secondary)' }}>
          Curated news articles from various RSS feeds, updated weekly.
        </p>
      </div>

      <DashboardBreadcrumb
        breadcrumb={breadcrumb}
        onNavigate={handleBreadcrumbClick}
        onRefreshFeed={(feedName) => {
          setCurrentPage(1);
          setHasMore(true);
          return loadArticlesInFeed(feedName, 1, false);
        }}
        setBreadcrumb={setBreadcrumb}
        setCurrentItems={setCurrentItems}
        setArticles={setArticles}
        loadArticlesInFeed={loadArticlesInFeed}
      />

      {/* Feeds Grid View */}
      {!isArticlesView && currentItems.length > 0 && (
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
          gap: '1.5rem',
          marginTop: '2rem'
        }}>
          {currentItems.map((item) => (
            <RSS
              key={item.name}
              item={item}
              handleItemClick={handleItemClick}
              getFeedCount={getFeedCount}
            />
          ))}
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
        <div style={{ marginTop: '2rem' }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: '3rem', color: 'var(--text-secondary)' }}>
              Loading articles...
            </div>
          ) : articles.length > 0 ? (
            <>
              <div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
                gap: '1.5rem',
                alignItems: 'stretch'
              }}>
                {articles.map((article) => (
                  <div key={article.id} style={{ height: '100%' }}>
                    <Article
                      article={article}
                      onView={() => handleViewArticle(article)}
                      isAdmin={isAdmin}
                      onDelete={() => handleDeleteArticle(article.id)}
                    />
                  </div>
                ))}
              </div>

              {/* Load More Button */}
              {hasMore && (
                <div style={{ display: 'flex', justifyContent: 'center', marginTop: '2rem' }}>
                  <button
                    onClick={handleLoadMore}
                    disabled={loadingMore}
                    style={{
                      backgroundColor: loadingMore ? 'var(--nav-active-bg-dim)' : 'var(--nav-active-bg)',
                      color: loadingMore ? 'var(--text-secondary)' : 'var(--nav-active-text)',
                      border: 'none',
                      cursor: loadingMore ? 'not-allowed' : 'pointer',
                      opacity: loadingMore ? 0.6 : 1
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
                        <div style={{
                          width: '16px',
                          height: '16px',
                          border: '2px solid currentColor',
                          borderTopColor: 'transparent',
                          borderRadius: '50%',
                          animation: 'spin 1s linear infinite'
                        }} />
                        Loading...
                      </>
                    ) : (
                      <>
                        <ChevronDown size={16} />
                        Load More
                      </>
                    )}
                  </button>
                </div>
              )}
            </>
          ) : (
            <div style={{
              textAlign: 'center',
              padding: '3rem',
              color: 'var(--text-secondary)',
              backgroundColor: 'var(--card-bg)',
              borderRadius: '12px',
              border: '1px solid var(--border)'
            }}>
              No articles found for this week.
            </div>
          )}
        </div>
      )}

      {/* Empty State */}
      {!isArticlesView && currentItems.length === 0 && (
        <div style={{
          textAlign: 'center',
          padding: '3rem',
          color: 'var(--text-secondary)',
          backgroundColor: 'var(--card-bg)',
          borderRadius: '12px',
          border: '1px solid var(--border)',
          marginTop: '2rem'
        }}>
          No items available.
        </div>
      )}
    </div>
  );
};

export default Home;
