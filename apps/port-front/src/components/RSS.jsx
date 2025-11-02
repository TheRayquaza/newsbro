import { useEffect, useState } from "react";
import { Link, FileText, Folder } from "lucide-react";

const RSS = ({
  item,
  handleItemClick,
  getFeedCount
}) => {
  const [feedCount, setFeedCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const isMetaFeed = item.link === "";
  const isFeed = item.link && item.link !== "";

  useEffect(() => {
    const loadCounts = async () => {
      setLoading(true);
      if (isMetaFeed) {
        const feeds = getFeedCount(item);
        setFeedCount(feeds);
      }
      setLoading(false);
    };
    loadCounts();
  }, [item, getFeedCount, isMetaFeed]);

  return (
    <button
      key={item.name}
      onClick={() => handleItemClick(item)}
      style={{
        backgroundColor: 'var(--nav-bg)',
        borderColor: 'var(--nav-border)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.borderColor = 'var(--nav-active-text)';
        e.currentTarget.querySelector('.hover-overlay').style.opacity = '1';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.borderColor = 'var(--nav-border)';
        e.currentTarget.querySelector('.hover-overlay').style.opacity = '0';
      }}
      className="group relative overflow-hidden rounded-2xl p-8 transition-all duration-300 text-left w-full"
    >
      <div 
        className="hover-overlay absolute inset-0 transition-opacity"
        style={{
          background: 'linear-gradient(to bottom right, var(--nav-active-bg), var(--search-button-active-bg))',
          opacity: '0'
        }}
      />
      <div className="relative z-10">
        <h3 
          style={{ color: 'var(--nav-active-text)' }}
          className="text-2xl font-bold transition mb-2"
        >
          {item.display_name}
        </h3>
        <p 
          style={{ color: 'var(--nav-text-muted)' }}
          className="transition mb-6 text-sm line-clamp-2"
        >
          {item.description}
        </p>
        <div className="flex items-center justify-between flex-wrap gap-3">
          <div className="flex items-center gap-3">
            {loading ? (
              <div 
                style={{
                  backgroundColor: 'var(--nav-active-bg)',
                  color: 'var(--nav-active-text)'
                }}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium"
              >
                <div 
                  style={{ borderTopColor: 'var(--nav-active-text)' }}
                  className="animate-spin rounded-full h-3 w-3 border-t-2"
                ></div>
                Loading...
              </div>
            ) : (
              <>
                {isMetaFeed && (
                  <div 
                    style={{
                      backgroundColor: 'var(--nav-active-bg)',
                      color: 'var(--nav-active-text)'
                    }}
                    className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition"
                  >
                    <Link className="w-4 h-4" />
                    {feedCount} {feedCount === 1 ? 'feed' : 'feeds'}
                  </div>
                )}
                {!isMetaFeed && !isFeed && item.children && (
                  <div 
                    style={{
                      backgroundColor: 'var(--nav-active-bg)',
                      color: 'var(--nav-active-text)'
                    }}
                    className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition"
                  >
                    <Folder className="w-4 h-4" />
                    {item.children.length} {item.children.length === 1 ? 'item' : 'items'}
                  </div>
                )}
              </>
            )}
          </div>
          <div>
            {isFeed ? (
              <span 
                style={{
                  backgroundColor: 'var(--success)',
                  color: '#ffffff',
                  borderColor: 'var(--success-hover)',
                  borderWidth: '1px',
                  borderStyle: 'solid'
                }}
                className="px-3 py-1 text-xs rounded-full"
              >
                RSS
              </span>
            ) : (
              <span 
                style={{
                  backgroundColor: 'var(--color-purple-500)',
                  color: '#ffffff',
                  borderColor: 'var(--nav-border)',
                  borderWidth: '1px',
                  borderStyle: 'solid'
                }}
                className="px-3 py-1 text-xs rounded-full"
              >
                Group
              </span>
            )}
          </div>
        </div>
      </div>
    </button>
  );
};

export default RSS;
