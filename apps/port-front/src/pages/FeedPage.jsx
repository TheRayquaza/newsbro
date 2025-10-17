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

    // Detect mobile/desktop
    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth < 768);
        };
        checkMobile();
        window.addEventListener('resize', checkMobile);
        return () => window.removeEventListener('resize', checkMobile);
    }, []);

    // Mock articles data
    useEffect(() => {
        const mockArticles = [
            {
                id: 1,
                title: "AI Breakthrough: New Language Model Achieves Human-Level Reasoning",
                abstract: "Researchers announce a significant advancement in artificial intelligence with a model that demonstrates unprecedented reasoning capabilities across multiple domains.",
                category: "Technology",
                subcategory: "AI & ML",
                published_at: "2025-10-15T10:30:00Z",
                link: "https://example.com/ai-breakthrough"
            },
            {
                id: 2,
                title: "Global Climate Summit Reaches Historic Agreement",
                abstract: "World leaders unite to commit to aggressive carbon reduction targets, marking a turning point in international climate policy and cooperation.",
                category: "Environment",
                subcategory: "Climate Change",
                published_at: "2025-10-14T14:20:00Z",
                link: "https://example.com/climate-summit"
            },
            {
                id: 3,
                title: "Quantum Computing Startup Secures $500M in Funding",
                abstract: "Silicon Valley startup announces massive funding round to scale quantum computing technology for commercial applications in healthcare and finance.",
                category: "Business",
                subcategory: "Startups",
                published_at: "2025-10-13T09:15:00Z",
                link: "https://example.com/quantum-funding"
            },
            {
                id: 4,
                title: "Scientists Discover Potential Treatment for Alzheimer's Disease",
                abstract: "Medical researchers report promising results from clinical trials of a new drug that shows significant improvements in cognitive function for patients.",
                category: "Health",
                subcategory: "Medical Research",
                published_at: "2025-10-12T16:45:00Z",
                link: "https://example.com/alzheimers-treatment"
            },
            {
                id: 5,
                title: "Electric Vehicle Sales Surpass Traditional Cars in Europe",
                abstract: "For the first time in history, electric vehicle sales have overtaken gasoline-powered cars across major European markets, signaling a major shift.",
                category: "Automotive",
                subcategory: "Electric Vehicles",
                published_at: "2025-10-11T11:30:00Z",
                link: "https://example.com/ev-sales"
            },
            {
                id: 6,
                title: "New Space Telescope Captures Never-Before-Seen Galaxies",
                abstract: "NASA's latest space observatory reveals stunning images of distant galaxies formed just after the Big Bang, offering insights into the early universe.",
                category: "Science",
                subcategory: "Space",
                published_at: "2025-10-10T08:00:00Z",
                link: "https://example.com/space-telescope"
            },
            {
                id: 7,
                title: "Breakthrough in Fusion Energy Brings Clean Power Closer",
                abstract: "Scientists achieve net energy gain in fusion reaction for the third consecutive time, demonstrating the viability of fusion as a clean energy source.",
                category: "Energy",
                subcategory: "Fusion",
                published_at: "2025-10-09T13:20:00Z",
                link: "https://example.com/fusion-energy"
            },
            {
                id: 8,
                title: "Major Tech Companies Announce AI Safety Coalition",
                abstract: "Leading technology firms form unprecedented alliance to establish ethical guidelines and safety standards for artificial intelligence development.",
                category: "Technology",
                subcategory: "AI Ethics",
                published_at: "2025-10-08T15:40:00Z",
                link: "https://example.com/ai-safety"
            }
        ];

        setArticles(mockArticles);
    }, []);

    // Keyboard navigation for desktop
    useEffect(() => {
        if (isMobile) return;

        const handleKeyDown = (e) => {
            if (removing) return;

            if (e.key === 'ArrowLeft' && currentIndex > 0) {
                setCurrentIndex(currentIndex - 1);
            } else if (e.key === 'ArrowRight' && currentIndex < articles.length - 1) {
                setCurrentIndex(currentIndex + 1);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [currentIndex, articles.length, removing, isMobile]);

    const handleSwipe = (direction) => {
        if (removing) return;
        if (direction !== 'left' && direction !== 'right') return;

        setRemoving(true);

        if (cardRef.current && isMobile) {
            const moveX = direction === 'right' ? 1000 : -1000;
            cardRef.current.style.transition = 'transform 0.3s ease-out';
            cardRef.current.style.transform = `translateX(${moveX}px) rotate(${direction === 'right' ? 20 : -20}deg)`;
        }

        setTimeout(() => {
            setArticles(prev => prev.filter((_, idx) => idx !== currentIndex));
            setRemoving(false);
            setDragOffset({ x: 0, y: 0 });

            if (cardRef.current) {
                cardRef.current.style.transition = '';
                cardRef.current.style.transform = '';
            }
        }, isMobile ? 300 : 0);

        if (currentIndex === articles.length - 1) {
            setCurrentIndex(prev => Math.max(0, prev - 1));
        }
    };

    const handlePrevious = () => {
        if (currentIndex > 0 && !removing) {
            setCurrentIndex(currentIndex - 1);
        }
    };

    const handleNext = () => {
        if (currentIndex < articles.length - 1 && !removing) {
            setCurrentIndex(currentIndex + 1);
        }
    };

    // Touch/Mouse handlers for mobile
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

        const threshold = 100;

        if (Math.abs(dragOffset.x) > threshold) {
            if (dragOffset.x > 0) {
                handleSwipe('right');
            } else {
                handleSwipe('left');
            }
        } else {
            setDragOffset({ x: 0, y: 0 });
        }
    };

    const rotation = isDragging && isMobile ? dragOffset.x / 20 : 0;
    const opacity = isDragging && isMobile ? Math.max(0.5, 1 - Math.abs(dragOffset.x) / 300) : 1;

    if (articles.length === 0) {
        return (
            <div className="flex-grow flex items-center justify-center p-4">
                <div className="text-center">
                    <div className="mb-6 flex justify-center">
                        <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-400 to-cyan-400 flex items-center justify-center">
                            <Heart className="w-12 h-12 text-white" fill="white" />
                        </div>
                    </div>
                    <h2 className="text-3xl font-bold text-white mb-4">All Caught Up!</h2>
                    <p className="text-slate-400">You've reviewed all articles in your feed.</p>
                </div>
            </div>
        );
    }

    return (
        <div ref={containerRef} className="flex items-center justify-center p-4 md:p-8 overflow-hidden">
            <div className="w-full max-w-6xl">
                {/* Desktop Carousel View */}
                {!isMobile && (
                    <div className="relative">
                        {/* Header */}
                        <div className="mb-8 text-center">
                            <div className="inline-flex items-center gap-4 px-6 py-3 bg-white/10 backdrop-blur-lg rounded-full border border-white/20">
                                <span className="text-white font-medium">
                                    {currentIndex + 1} / {articles.length}
                                </span>
                            </div>
                        </div>

                        {/* Carousel Container */}
                        <div className="relative flex items-center gap-6">
                            {/* Left Arrow */}
                            <button
                                onClick={handlePrevious}
                                disabled={currentIndex === 0 || removing}
                                className="flex-shrink-0 w-14 h-14 bg-white/10 hover:bg-white/20 backdrop-blur-lg rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100 border border-white/20 flex items-center justify-center group"
                            >
                                <ChevronLeft className="w-7 h-7 text-white group-hover:text-purple-300 transition-colors" strokeWidth={2.5} />
                            </button>

                            {/* Main Card */}
                            <div className="flex-1 transition-all duration-500 ease-out">
                                <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-3xl overflow-hidden border border-white/10 shadow-2xl">
                                    <div className="p-12">
                                        <h2 className="text-4xl font-bold text-white mb-6 leading-tight">
                                            {articles[currentIndex].title}
                                        </h2>

                                        <p className="text-slate-300 text-xl mb-8 leading-relaxed">
                                            {articles[currentIndex].abstract}
                                        </p>

                                        <div className="flex flex-wrap gap-3 mb-6">
                                            <span className="px-6 py-2 bg-gradient-to-r from-blue-400/20 to-cyan-400/20 border border-purple-500/30 text-purple-300 rounded-full text-sm font-semibold">
                                                {articles[currentIndex].category}
                                            </span>
                                            {articles[currentIndex].subcategory && (
                                                <span className="px-6 py-2 bg-gradient-to-r from-cyan-500/20 to-blue-500/20 border border-cyan-500/30 text-cyan-300 rounded-full text-sm font-semibold">
                                                    {articles[currentIndex].subcategory}
                                                </span>
                                            )}
                                        </div>

                                        {articles[currentIndex].published_at && (
                                            <div className="flex items-center gap-2 text-slate-400 mb-8">
                                                <Calendar className="w-5 h-5" />
                                                <span className="font-medium">
                                                    {new Date(articles[currentIndex].published_at).toLocaleDateString('en-US', {
                                                        year: 'numeric',
                                                        month: 'long',
                                                        day: 'numeric'
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
                                                    className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-cyan-400 hover:from-blue-500 hover:to-cyan-500 transition-all hover:scale-105 hover:underline underline-offset-4"
                                                >
                                                    <ExternalLink className="w-5 h-5 text-blue-400" />
                                                    <span>Read Full Article</span>
                                                </a>

                                            )}

                                            <button
                                                onClick={() => handleSwipe('right')}
                                                disabled={removing}
                                                className="inline-flex items-center gap-2 text-white px-8 py-3 bg-gradient-to-r from-green-400 to-emerald-500 hover:from-green-500 hover:to-emerald-600 rounded-full font-semibold transition-all shadow-lg hover:shadow-xl hover:scale-105 disabled:opacity-50"
                                            >
                                                <Heart className="w-5 h-5" fill="white" strokeWidth={0} />
                                                <span>Like</span>
                                            </button>

                                            <button
                                                onClick={() => handleSwipe('left')}
                                                disabled={removing}
                                                className="inline-flex items-center gap-2 text-white px-8 py-3 bg-gradient-to-r from-red-500 to-red-500 hover:from-red-600 hover:to-red-600 rounded-full font-semibold transition-all shadow-lg hover:shadow-xl hover:scale-105 disabled:opacity-50"
                                            >
                                                <X className="w-5 h-5" strokeWidth={3} />
                                                <span>Pass</span>
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Right Arrow */}
                            <button
                                onClick={handleNext}
                                disabled={currentIndex === articles.length - 1 || removing}
                                className="flex-shrink-0 w-14 h-14 bg-white/10 hover:bg-white/20 backdrop-blur-lg rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100 border border-white/20 flex items-center justify-center group"
                            >
                                <ChevronRight className="w-7 h-7 text-white group-hover:text-purple-300 transition-colors" strokeWidth={2.5} />
                            </button>
                        </div>

                        {/* Dots Indicator */}
                        <div className="flex justify-center gap-2 mt-8">
                            {articles.map((_, idx) => (
                                <button
                                    key={idx}
                                    onClick={() => !removing && setCurrentIndex(idx)}
                                    className={`transition-all rounded-full ${idx === currentIndex
                                        ? 'w-8 h-2 bg-gradient-to-r from-blue-400 to-cyan-400'
                                        : 'w-2 h-2 bg-white/30 hover:bg-white/50'
                                        }`}
                                />
                            ))}
                        </div>
                    </div>
                )}

                {/* Mobile Swipe View */}
                {isMobile && (
                    <div className="flex flex-col items-center justify-center">
                        <div className="mb-6 px-4 py-2 bg-white/10 backdrop-blur-lg rounded-full">
                            <span className="text-white font-medium text-sm">
                                {currentIndex + 1} / {articles.length}
                            </span>
                        </div>

                        <div className="relative w-full max-w-md h-[600px]">
                            <div
                                ref={cardRef}
                                className="absolute inset-0 cursor-grab active:cursor-grabbing"
                                style={{
                                    transform: `translateX(${dragOffset.x}px) rotate(${rotation}deg)`,
                                    opacity: opacity,
                                    transition: isDragging ? 'none' : 'transform 0.3s ease-out, opacity 0.3s ease-out'
                                }}
                                onMouseDown={(e) => handleStart(e.clientX, e.clientY)}
                                onMouseMove={(e) => handleMove(e.clientX, e.clientY)}
                                onMouseUp={handleEnd}
                                onMouseLeave={handleEnd}
                                onTouchStart={(e) => handleStart(e.touches[0].clientX, e.touches[0].clientY)}
                                onTouchMove={(e) => handleMove(e.touches[0].clientX, e.touches[0].clientY)}
                                onTouchEnd={handleEnd}
                            >
                                <div className="h-full rounded-3xl overflow-hidden bg-gradient-to-br from-slate-800 to-slate-900 border border-white/10 shadow-2xl">
                                    <div className="h-full overflow-y-auto p-6 flex flex-col">
                                        {isDragging && Math.abs(dragOffset.x) > 30 && (
                                            <div className="absolute top-8 left-0 right-0 flex justify-between px-8 pointer-events-none z-10">
                                                <div
                                                    className={`transform rotate-12 px-6 py-3 rounded-2xl border-4 text-2xl font-black backdrop-blur-sm ${dragOffset.x > 0
                                                        ? 'border-green-400 text-green-400 bg-green-500/20 scale-110'
                                                        : 'border-green-500/10 text-green-400/10'
                                                        }`}
                                                >
                                                    LIKE
                                                </div>
                                                <div
                                                    className={`transform -rotate-12 px-6 py-3 rounded-2xl border-4 text-2xl font-black backdrop-blur-sm ${dragOffset.x < 0
                                                        ? 'border-red-400 text-red-400 bg-red-500/20 scale-110'
                                                        : 'border-red-500/10 text-red-400/10'
                                                        }`}
                                                >
                                                    NOPE
                                                </div>
                                            </div>
                                        )}

                                        <div className="flex-1 flex flex-col justify-center">
                                            <h2 className="text-2xl font-bold text-white mb-4 leading-tight">
                                                {articles[currentIndex].title}
                                            </h2>

                                            <p className="text-slate-300 text-base mb-6 leading-relaxed">
                                                {articles[currentIndex].abstract}
                                            </p>

                                            <div className="flex flex-wrap gap-2 mb-4">
                                                <span className="px-4 py-1.5 bg-gradient-to-r from-blue-400/20 to-cyan-400/20 border border-purple-500/30 text-purple-300 rounded-full text-xs font-semibold">
                                                    {articles[currentIndex].category}
                                                </span>
                                                {articles[currentIndex].subcategory && (
                                                    <span className="px-4 py-1.5 bg-gradient-to-r from-cyan-500/20 to-blue-500/20 border border-cyan-500/30 text-cyan-300 rounded-full text-xs font-semibold">
                                                        {articles[currentIndex].subcategory}
                                                    </span>
                                                )}
                                            </div>

                                            {articles[currentIndex].published_at && (
                                                <div className="flex items-center gap-2 text-slate-400 mb-4">
                                                    <Calendar className="w-4 h-4" />
                                                    <span className="text-sm font-medium">
                                                        {new Date(articles[currentIndex].published_at).toLocaleDateString('en-US', {
                                                            year: 'numeric',
                                                            month: 'long',
                                                            day: 'numeric'
                                                        })}
                                                    </span>
                                                </div>
                                            )}

                                            {articles[currentIndex].link && (
                                                <a
                                                    href={articles[currentIndex].link}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="inline-flex items-center justify-center gap-2 text-white px-6 py-3  text-transparent bg-clip-text rounded-full font-semibold transition-all shadow-lg"
                                                    onClick={(e) => e.stopPropagation()}
                                                >
                                                    <ExternalLink className="w-4 h-4" />
                                                    <span className='text-sm text-center'>Read Full Article</span>
                                                </a>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="mt-8 flex justify-center items-center gap-6">
                            <button
                                onClick={() => handleSwipe('left')}
                                disabled={removing}
                                className="w-16 h-16 bg-gradient-to-r from-red-400 to-red-500 hover:from-red-400 hover:to-red-600 text-white rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-50 shadow-2xl flex items-center justify-center"
                            >
                                <X className="w-8 h-8" strokeWidth={3} />
                            </button>

                            <button
                                onClick={() => handleSwipe('right')}
                                disabled={removing}
                                className="w-16 h-16 bg-gradient-to-br from-green-400 to-emerald-500 hover:from-green-500 hover:to-emerald-600 text-white rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-50 shadow-2xl flex items-center justify-center"
                            >
                                <Heart className="w-8 h-8" fill="white" strokeWidth={0} />
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default FeedPage;
