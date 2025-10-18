import React, { useContext, useEffect, useState, useRef } from 'react';
import { Link } from 'lucide-react';
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

  // Data State
  const [articleFeedback, setArticleFeedback] = useState(null);
  const [selectedArticle, setSelectedArticle] = useState(null);
  const [currentItems, setCurrentItems] = useState([]);
  const [articles, setArticles] = useState([]);

  // Caching State
  const [articlesCache, setArticlesCache] = useState({});
  const [feedbackCache, setFeedbackCache] = useState({});
  const rssTreeCache = useRef(null);

  const pendingRequests = useRef({
    counts: {},
    articles: {},
    feedback: {}
  });

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

  const loadArticlesInFeed = async (feedName) => {
    if (articlesCache[feedName]) {
      setArticles(articlesCache[feedName]);
      return { articles: articlesCache[feedName] };
    }

    if (pendingRequests.current.articles[feedName]) {
      return pendingRequests.current.articles[feedName];
    }

    setLoading(true);

    const startOfWeek = getStartOfWeek();

    const requestPromise = api.getArticles({
      "feed_name": feedName,
      "start_date": startOfWeek.toISOString()
    })
      .then(response => {
        setArticles(response.articles);
        setArticlesCache(prev => ({ ...prev, [feedName]: response.articles }));
        setLoading(false);
        delete pendingRequests.current.articles[feedName];
        return response;
      })
      .catch(error => {
        alert("Failed to load articles for the selected feed.: " + error.message);
        setLoading(false);
        delete pendingRequests.current.articles[feedName];
        return { articles: [] };
      });

    pendingRequests.current.articles[feedName] = requestPromise;
    return requestPromise;
  };

  const handleItemClick = async (item) => {
    setLoading(true);

    if (item.children && item.children.length > 0) {
      setCurrentItems(item.children);
      setBreadcrumb([...breadcrumb, item]);
      setArticles([]);
      setLoading(false);
    }
    else if (item.link) {
      await loadArticlesInFeed(item.name);
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
        setArticlesCache(prev => ({ ...prev, [currentFeedName]: updatedArticles }));
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
              <h1 className="text-3xl font-bold bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
                News Feed Hub
              </h1>
            </div>
            <h2 className="text-slate-400">
              Curated news articles from various RSS feeds, updated weekly.
            </h2>
            <DashboardBreadcrumb
              breadcrumb={breadcrumb}
              setBreadcrumb={setBreadcrumb}
              setCurrentItems={setCurrentItems}
              setArticles={setArticles}
              loadArticlesInFeed={(feedName) => loadArticlesInFeed(feedName)}
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
            <div className="bg-gradient-to-r from-blue-500/10 to-cyan-500/10 border border-blue-500/20 rounded-xl p-6">
              <h2 className="text-2xl font-bold text-blue-300 mb-2">{currentLevel?.display_name}</h2>
              <p className="text-slate-400 text-sm">
                {currentLevel?.description}
              </p>
              {currentLevel?.link && (
                <a
                  href={currentLevel.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-block mt-2 text-xs text-green-400 hover:underline"
                >
                  <Link className="inline-block align-middle mr-2" size={14} /> RSS Feed
                </a>
              )}
            </div>

            {/* Articles List */}
            <div className="space-y-4">
              {loading ? (
                <div className="flex items-center justify-center py-12">
                  <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-blue-400"></div>
                </div>
              ) : articles.length > 0 ? (
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

export default Home;
