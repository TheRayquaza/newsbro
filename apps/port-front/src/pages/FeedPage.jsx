import React, { useState, useEffect, useRef } from 'react';
import { Calendar, ExternalLink, ChevronLeft, ChevronRight, Heart, X } from 'lucide-react';
import api from "../api/api";

const FeedPage = () => {
    // Data States
    const [articles, setArticles] = useState([]);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [models, setModels] = useState([]);
    const [selectedModel, setSelectedModel] = useState('');

    // UI States
    const [isLoading, setLoading] = useState(false);
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

    useEffect(() => {
        if (!selectedModel) return;

        const loadFeed = async () => {
            setLoading(true);
            try {
                const feedData = await api.getFeed({ model: selectedModel, limit: 20 });
                setArticles(feedData.articles || []);
                setCurrentIndex(0);
            } catch (error) {
                console.error('Error loading feed:', error);
            } finally {
                setLoading(false);
            }
        };

        loadFeed()
    }, [selectedModel]);

    useEffect(() => {
        const loadModels = async () => {
            try {
                const modelsData = await api.getFeedModels();
                setModels(modelsData || []);
                if (modelsData && modelsData.length > 0) {
                    setSelectedModel(modelsData[0]);
                }
            } catch (error) {
                console.error('Error loading models:', error);
            }
        };

        loadModels();
    }, []);


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

    const handleSwipe = async (direction) => {
        if (removing) return;
        if (direction !== 'left' && direction !== 'right') return;

        const isPositive = direction === 'right' ? true : false;
        const currentArticle = articles[currentIndex];
        if (!currentArticle) return;

        setRemoving(true);

        try {
            if (cardRef.current && isMobile) {
                const moveX = direction === 'right' ? 1000 : -1000;
                cardRef.current.style.transition = 'transform 0.3s ease-out';
                cardRef.current.style.transform = `translateX(${moveX}px) rotate(${direction === 'right' ? 20 : -20}deg)`;
            }

            await api.createArticleFeedback(currentArticle.id, { value: isPositive });
            await api.removeArticleFromFeed(currentArticle.id);

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

        } catch (error) {
            console.error('Error processing swipe:', error);
            setRemoving(false);
        }

        if (articles.length == 1 && !removing) {
            setLoading(true);
            loadMoreArticles().then(() => setLoading(false));
        }
    };

    const handlePrevious = () => {
        if (currentIndex > 0 && !removing) {
            setCurrentIndex(currentIndex - 1);
        }
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

    const handleModelChange = (e) => {
        const newModel = e.target.value;
        setSelectedModel(newModel);
    };

    const handleNext = async () => {
        if (currentIndex < articles.length - 1 && !removing) {
            setCurrentIndex(currentIndex + 1);
        } else if (currentIndex >= articles.length - 1 && !removing) {
            console.log('Loading more articles...');
            await loadMoreArticles();
        }
    };

    const loadMoreArticles = async () => {
        if (removing || isLoading) return;

        setLoading(true);
        try {
            const feedData = await api.getFeed({ model: selectedModel, limit: 20, offset: articles.length });
            setArticles((prev) => [...prev, ...(feedData.articles || [])]);
        } catch (error) {
            console.error('Error loading more articles:', error);
        } finally {
            setLoading(false);
        }
    };

    const rotation = isDragging && isMobile ? dragOffset.x / 20 : 0;
    const opacity = isDragging && isMobile ? Math.max(0.5, 1 - Math.abs(dragOffset.x) / 300) : 1;

    if (isLoading) {
        return (
            <div className="flex-grow flex items-center justify-center p-4">
                <div className="text-center">
                    <div
                        style={{
                            background: 'linear-gradient(to bottom right, var(--nav-gradient-from), var(--nav-gradient-to))'
                        }}
                        className="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 animate-pulse"
                    >
                        <Calendar className="w-8 h-8 text-white" />
                    </div>
                    <p style={{ color: 'var(--nav-text-muted)' }} className="text-lg">
                        Loading your feed...
                    </p>
                </div>
            </div>
        )
    }

    if (articles.length === 0) {
        return (
            <div className="flex items-center justify-center p-4 md:p-8 overflow-hidden">
                <div className="w-full max-w-6xl">
                    {/* Model Selector */}
                    {models.length > 0 && (
                        <div className="mb-6 flex justify-center">
                            <div
                                style={{
                                    backgroundColor: 'var(--nav-hover-bg)',
                                    backdropFilter: 'blur(16px)',
                                    borderColor: 'var(--nav-border)',
                                    borderWidth: '1px',
                                    borderStyle: 'solid'
                                }}
                                className="inline-flex items-center gap-3 px-4 py-2 rounded-full"
                            >
                                <span
                                    style={{ color: 'var(--nav-text-muted)' }}
                                    className="text-sm font-medium"
                                >
                                    Model:
                                </span>
                                <select
                                    value={selectedModel}
                                    onChange={handleModelChange}
                                    disabled={removing}
                                    style={{
                                        backgroundColor: 'var(--nav-active-bg)',
                                        color: 'var(--nav-text)',
                                        borderColor: 'var(--nav-border)',
                                        borderWidth: '1px',
                                        borderStyle: 'solid'
                                    }}
                                    className="px-3 py-1.5 rounded-full text-sm font-medium cursor-pointer focus:outline-none focus:ring-2 focus:ring-offset-2"
                                    onFocus={(e) => e.target.style.borderColor = 'var(--nav-active-text)'}
                                    onBlur={(e) => e.target.style.borderColor = 'var(--nav-border)'}
                                >
                                    {models.map((model) => (
                                        <option key={model} value={model}>
                                            {model}
                                        </option>
                                    ))}
                                </select>
                            </div>
                        </div>
                    )}

                    {/* Centered Empty State */}
                    <div className="flex items-center justify-center" style={{ minHeight: '50vh' }}>
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
                </div>
            </div>
        );
    }

    return (
        <div ref={containerRef} className="flex items-center justify-center p-4 md:p-8 overflow-hidden">
            <div className="w-full max-w-6xl">
                {/* Model Selector */}
                {models.length > 0 && (
                    <div className="mb-6 flex justify-center">
                        <div
                            style={{
                                backgroundColor: 'var(--nav-hover-bg)',
                                backdropFilter: 'blur(16px)',
                                borderColor: 'var(--nav-border)',
                                borderWidth: '1px',
                                borderStyle: 'solid'
                            }}
                            className="inline-flex items-center gap-3 px-4 py-2 rounded-full"
                        >
                            <span
                                style={{ color: 'var(--nav-text-muted)' }}
                                className="text-sm font-medium"
                            >
                                Model:
                            </span>
                            <select
                                value={selectedModel}
                                onChange={handleModelChange}
                                disabled={removing}
                                style={{
                                    backgroundColor: 'var(--nav-active-bg)',
                                    color: 'var(--nav-text)',
                                    borderColor: 'var(--nav-border)',
                                    borderWidth: '1px',
                                    borderStyle: 'solid'
                                }}
                                className="px-3 py-1.5 rounded-full text-sm font-medium cursor-pointer focus:outline-none focus:ring-2 focus:ring-offset-2"
                                onFocus={(e) => e.target.style.borderColor = 'var(--nav-active-text)'}
                                onBlur={(e) => e.target.style.borderColor = 'var(--nav-border)'}
                            >
                                {models.map((model) => (
                                    <option key={model} value={model}>
                                        {model}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>
                )}

                {/* Desktop Carousel View */}
                {!isMobile && (
                    <div className="relative">
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
