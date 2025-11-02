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
                        <div 
                            style={{
                                background: 'linear-gradient(to bottom right, var(--nav-gradient-from), var(--nav-gradient-to))'
                            }}
                            className="w-20 h-20 rounded-full flex items-center justify-center"
                        >
                            <Heart className="w-12 h-12 text-white" fill="white" />
                        </div>
                    </div>
                    <h2 
                        style={{ color: 'var(--nav-text)' }}
                        className="text-3xl font-bold mb-4"
                    >
                        All Caught Up!
                    </h2>
                    <p style={{ color: 'var(--nav-text-muted)' }}>
                        You've reviewed all articles in your feed.
                    </p>
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
                            <div 
                                style={{
                                    backgroundColor: 'var(--nav-hover-bg)',
                                    backdropFilter: 'blur(16px)',
                                    borderColor: 'var(--nav-border)',
                                    borderWidth: '1px',
                                    borderStyle: 'solid'
                                }}
                                className="inline-flex items-center gap-4 px-6 py-3 rounded-full"
                            >
                                <span 
                                    style={{ color: 'var(--nav-text)' }}
                                    className="font-medium"
                                >
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
                                style={{
                                    backgroundColor: 'var(--nav-hover-bg)',
                                    backdropFilter: 'blur(16px)',
                                    borderColor: 'var(--nav-border)',
                                    borderWidth: '1px',
                                    borderStyle: 'solid'
                                }}
                                onMouseEnter={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = 'var(--nav-active-bg)')}
                                onMouseLeave={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)')}
                                className="flex-shrink-0 w-14 h-14 rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center group"
                            >
                                <ChevronLeft 
                                    style={{ color: 'var(--nav-text)' }}
                                    className="w-7 h-7 transition-colors" 
                                    strokeWidth={2.5} 
                                />
                            </button>

                            {/* Main Card */}
                            <div className="flex-1 transition-all duration-500 ease-out">
                                <div 
                                    style={{
                                        backgroundColor: 'var(--nav-bg)',
                                        borderColor: 'var(--nav-border)',
                                        borderWidth: '1px',
                                        borderStyle: 'solid'
                                    }}
                                    className="rounded-3xl overflow-hidden shadow-2xl"
                                >
                                    <div className="p-12">
                                        <h2 
                                            style={{ color: 'var(--nav-text)' }}
                                            className="text-4xl font-bold mb-6 leading-tight"
                                        >
                                            {articles[currentIndex].title}
                                        </h2>

                                        <p 
                                            style={{ color: 'var(--nav-text-muted)' }}
                                            className="text-xl mb-8 leading-relaxed"
                                        >
                                            {articles[currentIndex].abstract}
                                        </p>

                                        <div className="flex flex-wrap gap-3 mb-6">
                                            <span 
                                                style={{
                                                    backgroundColor: 'var(--nav-active-bg)',
                                                    borderColor: 'var(--nav-border)',
                                                    color: 'var(--nav-active-text)',
                                                    borderWidth: '1px',
                                                    borderStyle: 'solid'
                                                }}
                                                className="px-6 py-2 rounded-full text-sm font-semibold"
                                            >
                                                {articles[currentIndex].category}
                                            </span>
                                            {articles[currentIndex].subcategory && (
                                                <span 
                                                    style={{
                                                        backgroundColor: 'var(--search-button-active-bg)',
                                                        borderColor: 'var(--nav-border)',
                                                        color: 'var(--nav-active-text)',
                                                        borderWidth: '1px',
                                                        borderStyle: 'solid'
                                                    }}
                                                    className="px-6 py-2 rounded-full text-sm font-semibold"
                                                >
                                                    {articles[currentIndex].subcategory}
                                                </span>
                                            )}
                                        </div>

                                        {articles[currentIndex].published_at && (
                                            <div 
                                                style={{ color: 'var(--nav-text-muted)' }}
                                                className="flex items-center gap-2 mb-8"
                                            >
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
                                                    style={{ color: 'var(--nav-active-text)' }}
                                                    onMouseEnter={(e) => e.currentTarget.style.color = 'var(--nav-gradient-to)'}
                                                    onMouseLeave={(e) => e.currentTarget.style.color = 'var(--nav-active-text)'}
                                                    className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all hover:scale-105 hover:underline underline-offset-4"
                                                >
                                                    <ExternalLink className="w-5 h-5" />
                                                    <span>Read Full Article</span>
                                                </a>
                                            )}

                                            <button
                                                onClick={() => handleSwipe('right')}
                                                disabled={removing}
                                                style={{
                                                    backgroundColor: 'var(--success)',
                                                    color: '#ffffff'
                                                }}
                                                onMouseEnter={(e) => !removing && (e.currentTarget.style.backgroundColor = 'var(--success-hover)')}
                                                onMouseLeave={(e) => !removing && (e.currentTarget.style.backgroundColor = 'var(--success)')}
                                                className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all shadow-lg hover:shadow-xl hover:scale-105 disabled:opacity-50"
                                            >
                                                <Heart className="w-5 h-5" fill="white" strokeWidth={0} />
                                                <span>Like</span>
                                            </button>

                                            <button
                                                onClick={() => handleSwipe('left')}
                                                disabled={removing}
                                                style={{
                                                    backgroundColor: 'var(--error)',
                                                    color: '#ffffff'
                                                }}
                                                onMouseEnter={(e) => !removing && (e.currentTarget.style.backgroundColor = 'var(--error-hover)')}
                                                onMouseLeave={(e) => !removing && (e.currentTarget.style.backgroundColor = 'var(--error)')}
                                                className="inline-flex items-center gap-2 px-8 py-3 rounded-full font-semibold transition-all shadow-lg hover:shadow-xl hover:scale-105 disabled:opacity-50"
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
                                style={{
                                    backgroundColor: 'var(--nav-hover-bg)',
                                    backdropFilter: 'blur(16px)',
                                    borderColor: 'var(--nav-border)',
                                    borderWidth: '1px',
                                    borderStyle: 'solid'
                                }}
                                onMouseEnter={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = 'var(--nav-active-bg)')}
                                onMouseLeave={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = 'var(--nav-hover-bg)')}
                                className="flex-shrink-0 w-14 h-14 rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center group"
                            >
                                <ChevronRight 
                                    style={{ color: 'var(--nav-text)' }}
                                    className="w-7 h-7 transition-colors" 
                                    strokeWidth={2.5} 
                                />
                            </button>
                        </div>

                        {/* Dots Indicator */}
                        <div className="flex justify-center gap-2 mt-8">
                            {articles.map((_, idx) => (
                                <button
                                    key={idx}
                                    onClick={() => !removing && setCurrentIndex(idx)}
                                    style={{
                                        backgroundColor: idx === currentIndex 
                                            ? 'var(--nav-active-text)' 
                                            : 'var(--nav-text-muted)',
                                        opacity: idx === currentIndex ? 1 : 0.3
                                    }}
                                    onMouseEnter={(e) => idx !== currentIndex && (e.currentTarget.style.opacity = '0.5')}
                                    onMouseLeave={(e) => idx !== currentIndex && (e.currentTarget.style.opacity = '0.3')}
                                    className={`transition-all rounded-full ${idx === currentIndex ? 'w-8 h-2' : 'w-2 h-2'}`}
                                />
                            ))}
                        </div>
                    </div>
                )}

                {/* Mobile Swipe View */}
                {isMobile && (
                    <div className="flex flex-col items-center justify-center">
                        <div 
                            style={{
                                backgroundColor: 'var(--nav-hover-bg)',
                                backdropFilter: 'blur(16px)'
                            }}
                            className="mb-6 px-4 py-2 rounded-full"
                        >
                            <span 
                                style={{ color: 'var(--nav-text)' }}
                                className="font-medium text-sm"
                            >
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
                                <div 
                                    style={{
                                        backgroundColor: 'var(--nav-bg)',
                                        borderColor: 'var(--nav-border)',
                                        borderWidth: '1px',
                                        borderStyle: 'solid'
                                    }}
                                    className="h-full rounded-3xl overflow-hidden shadow-2xl"
                                >
                                    <div className="h-full overflow-y-auto p-6 flex flex-col">
                                        {isDragging && Math.abs(dragOffset.x) > 30 && (
                                            <div className="absolute top-8 left-0 right-0 flex justify-between px-8 pointer-events-none z-10">
                                                <div
                                                    style={{
                                                        borderColor: dragOffset.x > 0 ? 'var(--success)' : 'transparent',
                                                        color: dragOffset.x > 0 ? 'var(--success)' : 'var(--nav-text-muted)',
                                                        backgroundColor: dragOffset.x > 0 ? 'var(--success-bg-hover)' : 'transparent',
                                                        borderWidth: '4px',
                                                        borderStyle: 'solid',
                                                        backdropFilter: 'blur(8px)',
                                                        opacity: dragOffset.x > 0 ? 1 : 0.2
                                                    }}
                                                    className="transform rotate-12 px-6 py-3 rounded-2xl text-2xl font-black"
                                                >
                                                    LIKE
                                                </div>
                                                <div
                                                    style={{
                                                        borderColor: dragOffset.x < 0 ? 'var(--error)' : 'transparent',
                                                        color: dragOffset.x < 0 ? 'var(--error)' : 'var(--nav-text-muted)',
                                                        backgroundColor: dragOffset.x < 0 ? 'var(--error-bg-hover)' : 'transparent',
                                                        borderWidth: '4px',
                                                        borderStyle: 'solid',
                                                        backdropFilter: 'blur(8px)',
                                                        opacity: dragOffset.x < 0 ? 1 : 0.2
                                                    }}
                                                    className="transform -rotate-12 px-6 py-3 rounded-2xl text-2xl font-black"
                                                >
                                                    NOPE
                                                </div>
                                            </div>
                                        )}

                                        <div className="flex-1 flex flex-col justify-center">
                                            <h2 
                                                style={{ color: 'var(--nav-text)' }}
                                                className="text-2xl font-bold mb-4 leading-tight"
                                            >
                                                {articles[currentIndex].title}
                                            </h2>

                                            <p 
                                                style={{ color: 'var(--nav-text-muted)' }}
                                                className="text-base mb-6 leading-relaxed"
                                            >
                                                {articles[currentIndex].abstract}
                                            </p>

                                            <div className="flex flex-wrap gap-2 mb-4">
                                                <span 
                                                    style={{
                                                        backgroundColor: 'var(--nav-active-bg)',
                                                        borderColor: 'var(--nav-border)',
                                                        color: 'var(--nav-active-text)',
                                                        borderWidth: '1px',
                                                        borderStyle: 'solid'
                                                    }}
                                                    className="px-4 py-1.5 rounded-full text-xs font-semibold"
                                                >
                                                    {articles[currentIndex].category}
                                                </span>
                                                {articles[currentIndex].subcategory && (
                                                    <span 
                                                        style={{
                                                            backgroundColor: 'var(--search-button-active-bg)',
                                                            borderColor: 'var(--nav-border)',
                                                            color: 'var(--nav-active-text)',
                                                            borderWidth: '1px',
                                                            borderStyle: 'solid'
                                                        }}
                                                        className="px-4 py-1.5 rounded-full text-xs font-semibold"
                                                    >
                                                        {articles[currentIndex].subcategory}
                                                    </span>
                                                )}
                                            </div>

                                            {articles[currentIndex].published_at && (
                                                <div 
                                                    style={{ color: 'var(--nav-text-muted)' }}
                                                    className="flex items-center gap-2 mb-4"
                                                >
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
                                                    style={{ color: 'var(--nav-active-text)' }}
                                                    className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-full font-semibold transition-all shadow-lg"
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
                                style={{
                                    backgroundColor: 'var(--error)',
                                    color: '#ffffff'
                                }}
                                className="w-16 h-16 rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-50 shadow-2xl flex items-center justify-center"
                            >
                                <X className="w-8 h-8" strokeWidth={3} />
                            </button>

                            <button
                                onClick={() => handleSwipe('right')}
                                disabled={removing}
                                style={{
                                    backgroundColor: 'var(--success)',
                                    color: '#ffffff'
                                }}
                                className="w-16 h-16 rounded-full transition-all hover:scale-110 active:scale-95 disabled:opacity-50 shadow-2xl flex items-center justify-center"
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
