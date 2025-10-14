import { Home } from "lucide-react";

const DashboardBreadcrumb = ({ 
    breadcrumb, 
    setBreadcrumb, 
    setCurrentItems, 
    setArticles, 
    loadArticlesInFeed,
 }) => {
    return (
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
                                        setArticles(loadArticlesInFeed(b.name));
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
    );
};

export default DashboardBreadcrumb;
