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
      className="group relative overflow-hidden bg-gradient-to-br from-slate-800/40 to-slate-800/20 border border-blue-500/20 rounded-2xl p-8 hover:border-blue-500/50 transition-all duration-300 text-left w-full"
    >
      <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-cyan-500/10 opacity-0 group-hover:opacity-100 transition-opacity" />
      <div className="relative z-10">
        <h3 className="text-2xl font-bold text-blue-300 group-hover:text-blue-200 transition mb-2">
          {item.display_name}
        </h3>
        <p className="text-slate-400 group-hover:text-slate-300 transition mb-6 text-sm line-clamp-2">
          {item.description}
        </p>
        <div className="flex items-center justify-between flex-wrap gap-3">
          <div className="flex items-center gap-3">
            {loading ? (
              <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-500/20 text-blue-300 rounded-lg text-sm font-medium">
                <div className="animate-spin rounded-full h-3 w-3 border-t-2 border-blue-400"></div>
                Loading...
              </div>
            ) : (
              <>
                {isMetaFeed && (
                  <div className="inline-flex items-center gap-2 px-4 py-2 bg-purple-500/20 group-hover:bg-purple-500/30 text-purple-300 rounded-lg text-sm font-medium transition">
                    <Link className="w-4 h-4" />
                    {feedCount} {feedCount === 1 ? 'feed' : 'feeds'}
                  </div>
                )}
                {!isMetaFeed && !isFeed && item.children && (
                  <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-500/20 group-hover:bg-blue-500/30 text-blue-300 rounded-lg text-sm font-medium transition">
                    <Folder className="w-4 h-4" />
                    {item.children.length} {item.children.length === 1 ? 'item' : 'items'}
                  </div>
                )}
              </>
            )}
          </div>
          <div>
            {isFeed ? (
              <span className="px-3 py-1 text-xs rounded-full bg-green-500/20 text-green-300 border border-green-500/30">
                RSS
              </span>
            ) : (
              <span className="px-3 py-1 text-xs rounded-full bg-purple-500/20 text-purple-300 border border-purple-500/30">
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