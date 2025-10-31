import React, { useState, useEffect, useRef } from 'react';
import { Calendar, ExternalLink, ChevronLeft, ChevronRight, Heart, X } from 'lucide-react';

const FeedPage = () => {
    const [articles, setArticles] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });
    const [isDragging, setIsDragging] = useState(false);
    const [removing, setRemoving] = useState(false);
    const [isMobile, setIsMobile] = useState(false);
    const cardRef = useRef(null);
    const startPos = useRef({ x: 0, y: 0 });
    const containerRef = useRef(null);

    useEffect(() => {
        const checkMobile = () => setIsMobile(window.innerWidth < 768);
        checkMobile();
        window.addEventListener('resize', checkMobile);
        return () => window.removeEventListener('resize', checkMobile);
    }, []);

    useEffect(() => {
        const mockArticles = [
            {
                id: 1,
                title: "AI Breakthrough: New Language Model Achieves Human-Level Reasoning",
                abstract: "Researchers announce a significant advancement in artificial intelligence...",
                category: "Technology",
                subcategory: "AI & ML",
                published_at: "2025-10-15T10:30:00Z",
                link: "https://example.com/ai-breakthrough"
            },
            // ... other mock articles ...
        ];
        setArticles(mockArticles);
    }, []);

    useEffect(() => {
        if (isMobile) return;
        const handleKeyDown = (e) => {
            if (removing) return;
            if (e.key === 'ArrowLeft' && currentIndex > 0) setCurrentIndex(currentIndex - 1);
            else if (e.key === 'ArrowRight' && currentIndex < articles.length - 1)
                setCurrentIndex(currentIndex + 1);
        };
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [currentIndex, articles.length, removing, isMobile]);

    const handleSwipe = (direction) => {
        if (removing || (direction !== 'left' && direction !== 'right')) return;
        setRemoving(true);

        if (cardRef.current && isMobile) {
            const moveX = direction === 'right' ? 1000 : -1000;
            cardRef.current.style.transition = 'transform 0.3s ease-out';
            cardRef.current.style.transform = `translateX(${moveX}px) rotate(${direction === 'right' ? 20 : -20}deg)`;
        }

        setTimeout(() => {
            setArticles((prev) => prev.filter((_, idx) => idx !== currentIndex));
            setRemoving(false);
            setDragOffset({ x: 0, y: 0 });
            if (cardRef.current) {
                cardRef.current.style.transition = '';
                cardRef.current.style.transform = '';
            }
        }, isMobile ? 300 : 0);

        if (currentIndex === articles.length - 1) setCurrentIndex((prev) => Math.max(0, prev - 1));
    };

    const handlePrevious = () => {
        if (currentIndex > 0 && !removing) setCurrentIndex(currentIndex - 1);
    };

    const handleNext = () => {
        if (currentIndex < articles.length - 1 && !removing) setCurrentIndex(currentIndex + 1);
    };

    const handleStart = (clientX, clientY) => {
        if (removing || !isMobile) return;
        setIsDragging(true);
        startPos.current = { x: clientX, y: clientY };
    };

    const handleMove = (clientX, clientY) => {
        if (!isDragging || removing || !isMobile) return;
        const deltaX = clientX - startPos.current.x;
        const deltaY = clientY - startPos.current.y;
        setDragOffset({ x: deltaX, y: deltaY });
    };

    const handleEnd = () => {
        if (!isDragging || removing || !isMobile) return;
        setIsDragging(false);
        if (Math.abs(dragOffset.x) > 100) handleSwipe(dragOffset.x > 0 ? 'right' : 'left');
        else setDragOffset({ x: 0, y: 0 });
    };

    const rotation = isDragging && isMobile ? dragOffset.x / 20 : 0;
    const opacity = isDragging && isMobile ? Math.max(0.5, 1 - Math.abs(dragOffset.x) / 300) : 1;

    if (articles.length === 0) {
        return (
            <div className="flex-grow flex items-center justify-center p-4">
                <div className="text-center">
                    <div className="mb-6 flex justify-center">
                        <div
                            className="w-20 h-20 rounded-full flex items-center justify-center"
                            style={{
                                background: `linear-gradient(to bottom right, var(--nav-gradient-from), var(--nav-gradient-to))`,
                            }}
                        >
                            <Heart className="w-12 h-12" style={{ color: 'var(--text-primary)' }} fill="white" />
                        </div>
                    </div>
                    <h2 className="text-3xl font-bold mb-4" style={{ color: 'var(--text-primary)' }}>
                        All Caught Up!
                    </h2>
                    <p style={{ color: 'var(--text-secondary)' }}>
                        You've reviewed all articles in your feed.
                    </p>
                </div>
            </div>
        );
    }

    return (
        <div ref={containerRef} className="flex items-center justify-center p-4 md:p-8 overflow-hidden">
            <div className="w-full max-w-6xl">
                {!isMobile && (
                    <div className="relative">
                        <div className="mb-8 text-center">
                            <div
                                className="inline-flex items-center gap-4 px-6 py-3 rounded-full border backdrop-blur-lg"
                                style={{
                                    backgroundColor: 'var(--bg-secondary)',
                                    borderColor: 'var(--nav-border)',
                                }}
                            >
                                <span className="font-medium" style={{ color: 'var(--text-primary)' }}>
                                    {currentIndex + 1} / {articles.length}
                                </span>
                            </div>
                        </div>

                        <div className="relative flex items-center gap-6">
                            <button
                                onClick={handlePrevious}
                                disabled={currentIndex === 0 || removing}
                                className="flex-shrink-0 w-14 h-14 rounded-full transition-all flex items-center justify-center group border backdrop-blur-lg"
                                style={{
                                    backgroundColor: 'var(--bg-secondary)',
                                    borderColor: 'var(--nav-border)',
                                }}
                            >
                                <ChevronLeft
                                    className="w-7 h-7 transition-colors"
                                    style={{ color: 'var(--text-primary)' }}
                                    strokeWidth={2.5}
                                />
                            </button>

                            <div className="flex-1 transition-all duration-500 ease-out">
                                <div
                                    className="rounded-3xl overflow-hidden border shadow-2xl"
                                    style={{
                                        backgroundColor: 'var(--bg-primary)',
                                        borderColor: 'var(--nav-border)',
                                    }}
                                >
                                    <div className="p-12">
                                        <h2
                                            className="text-4xl font-bold mb-6 leading-tight"
                                            style={{ color: 'var(--text-primary)' }}
                                        >
                                            {articles[currentIndex].title}
                                        </h2>

                                        <p
                                            className="text-xl mb-8 leading-relaxed"
                                            style={{ color: 'var(--text-secondary)' }}
                                        >
                                            {articles[currentIndex].abstract}
                                        </p>

                                        <div className="flex flex-wrap gap-3 mb-6">
                                            <span
                                                className="px-6 py-2 rounded-full text-sm font-semibold border"
                                                style={{
                                                    background: 'linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))',
                                                    borderColor: 'var(--nav-border)',
                                                    color: 'var(--text-primary)',
                                                }}
                                            >
                                                {articles[currentIndex].category}
                                            </span>
                                            {articles[currentIndex].subcategory && (
                                                <span
                                                    className="px-6 py-2 rounded-full text-sm font-semibold border"
                                                    style={{
                                                        background: 'linear-gradient(to right, var(--nav-gradient-to), var(--nav-gradient-from))',
                                                        borderColor: 'var(--nav-border)',
                                                        color: 'var(--text-primary)',
                                                    }}
                                                >
                                                    {articles[currentIndex].subcategory}
                                                </span>
                                            )}
                                        </div>

                                        {articles[currentIndex].published_at && (
                                            <div className="flex items-center gap-2 mb-8" style={{ color: 'var(--text-secondary)' }}>
                                                <Calendar className="w-5 h-5" />
                                                <span className="font-medium">
                                                    {new Date(articles[currentIndex].published_at).toLocaleDateString('en-US', {
                                                        year: 'numeric',
                                                        month: 'long',
                                                        day: 'numeric',
                                                    })}
                                                </span>
                                            </div>
                                        )}

                                        <div className="flex gap-4">
                                            {articles[currentIndex].link && (
                                                <a
                                                    href={articles[currentIndex].link}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all hover:scale-105"
                                                    style={{
                                                        color: 'var(--text-primary)',
                                                        background: 'linear-gradient(to right, var(--nav-gradient-from), var(--nav-gradient-to))',
                                                    }}
                                                >
                                                    <ExternalLink className="w-5 h-5" />
                                                    <span>Read Full Article</span>
                                                </a>
                                            )}

                                            <button
                                                onClick={() => handleSwipe('right')}
                                                disabled={removing}
                                                className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all hover:scale-105 shadow-lg disabled:opacity-50"
                                                style={{
                                                    color: 'white',
                                                    background: 'linear-gradient(to right, var(--success), var(--success-hover))',
                                                }}
                                            >
                                                <Heart className="w-5 h-5" fill="white" strokeWidth={0} />
                                                <span>Like</span>
                                            </button>

                                            <button
                                                onClick={() => handleSwipe('left')}
                                                disabled={removing}
                                                className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all hover:scale-105 shadow-lg disabled:opacity-50"
                                                style={{
                                                    color: 'white',
                                                    background: 'linear-gradient(to right, var(--error), var(--error-hover))',
                                                }}
                                            >
                                                <X className="w-5 h-5" strokeWidth={3} />
                                                <span>Pass</span>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <button
                                onClick={handleNext}
                                disabled={currentIndex === articles.length - 1 || removing}
                                className="flex-shrink-0 w-14 h-14 rounded-full transition-all flex items-center justify-center group border backdrop-blur-lg"
                                style={{
                                    backgroundColor: 'var(--bg-secondary)',
                                    borderColor: 'var(--nav-border)',
                                }}
                            >
                                <ChevronRight
                                    className="w-7 h-7 transition-colors"
                                    style={{ color: 'var(--text-primary)' }}
                                    strokeWidth={2.5}
                                />
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default FeedPage;
