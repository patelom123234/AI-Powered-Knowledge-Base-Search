import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';

import RelevantArticles from './RelevantArticles';

describe('RelevantArticles', () => {
    
  it('should render a grid of article cards from props', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'First Test Article' },
      { id: 'kb-002', title: 'Second Test Article' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('First Test Article')).toBeInTheDocument();
    expect(screen.getByText('Second Test Article')).toBeInTheDocument();
  });

  it('should render the articles header with icon and title', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'Test Article' }
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    const articlesIcon = screen.getByText('📚');
    const articlesTitle = screen.getByText('Relevant Articles');
    
    expect(articlesIcon).toBeInTheDocument();
    expect(articlesTitle).toBeInTheDocument();
  });

  it('should display article count in header', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'First Article' },
      { id: 'kb-002', title: 'Second Article' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('2 articles')).toBeInTheDocument();
  });

  it('should display singular "article" for single article', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'Single Article' }
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('1 article')).toBeInTheDocument();
  });

  it('should render article cards with proper structure', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'Test Article' }
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    // Check article number
    expect(screen.getByText('#1')).toBeInTheDocument();
    
    // Check article ID
    expect(screen.getByText('kb-001')).toBeInTheDocument();
    
    // Check article title
    expect(screen.getByText('Test Article')).toBeInTheDocument();
    
    // Check knowledge base badge
    expect(screen.getByText('📄')).toBeInTheDocument();
    expect(screen.getByText('Knowledge Base')).toBeInTheDocument();
  });

  it('should render multiple article cards with correct numbering', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'First Article' },
      { id: 'kb-002', title: 'Second Article' },
      { id: 'kb-003', title: 'Third Article' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('#1')).toBeInTheDocument();
    expect(screen.getByText('#2')).toBeInTheDocument();
    expect(screen.getByText('#3')).toBeInTheDocument();
    
    expect(screen.getByText('kb-001')).toBeInTheDocument();
    expect(screen.getByText('kb-002')).toBeInTheDocument();
    expect(screen.getByText('kb-003')).toBeInTheDocument();
  });

  it('should display a proper empty state when no articles are provided', () => {
    const mockArticles = [];
    
    render(<RelevantArticles articles={mockArticles} />);

    // Check header elements are still present
    expect(screen.getByText('📚')).toBeInTheDocument();
    expect(screen.getByText('Relevant Articles')).toBeInTheDocument();
    
    // Check empty state elements
    expect(screen.getByText('🔍')).toBeInTheDocument();
    expect(screen.getByText('No relevant articles found for this query.')).toBeInTheDocument();
  });

  it('should display empty state when articles is null', () => {
    render(<RelevantArticles articles={null} />);

    expect(screen.getByText('No relevant articles found for this query.')).toBeInTheDocument();
  });

  it('should display empty state when articles is undefined', () => {
    render(<RelevantArticles articles={undefined} />);

    expect(screen.getByText('No relevant articles found for this query.')).toBeInTheDocument();
  });

  it('should have proper CSS classes for styling', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'Test Article' }
    ];

    const { container } = render(<RelevantArticles articles={mockArticles} />);

    // Check for main container class
    const relevantArticles = container.querySelector('.relevant-articles');
    expect(relevantArticles).toBeInTheDocument();

    // Check for header classes
    const articlesHeader = container.querySelector('.articles-header');
    expect(articlesHeader).toBeInTheDocument();

    // Check for grid class
    const articlesGrid = container.querySelector('.articles-grid');
    expect(articlesGrid).toBeInTheDocument();

    // Check for card classes
    const articleCard = container.querySelector('.article-card');
    expect(articleCard).toBeInTheDocument();

    // Check for article header classes
    const articleHeader = container.querySelector('.article-header');
    expect(articleHeader).toBeInTheDocument();

    // Check for article footer classes
    const articleFooter = container.querySelector('.article-footer');
    expect(articleFooter).toBeInTheDocument();

    // Check for badge classes
    const articleBadge = container.querySelector('.article-badge');
    expect(articleBadge).toBeInTheDocument();
  });

  it('should have proper CSS classes for empty state', () => {
    const mockArticles = [];

    const { container } = render(<RelevantArticles articles={mockArticles} />);

    // Check for main container class
    const relevantArticles = container.querySelector('.relevant-articles');
    expect(relevantArticles).toBeInTheDocument();

    // Check for header classes
    const articlesHeader = container.querySelector('.articles-header');
    expect(articlesHeader).toBeInTheDocument();

    // Check for empty state classes
    const noArticles = container.querySelector('.no-articles');
    expect(noArticles).toBeInTheDocument();
  });

  it('should handle articles with different ID formats', () => {
    const mockArticles = [
      { id: 'article-123', title: 'Article with different ID format' },
      { id: 'KB_456', title: 'Another article with different ID' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('article-123')).toBeInTheDocument();
    expect(screen.getByText('KB_456')).toBeInTheDocument();
    expect(screen.getByText('Article with different ID format')).toBeInTheDocument();
    expect(screen.getByText('Another article with different ID')).toBeInTheDocument();
  });

  it('should handle articles with long titles', () => {
    const longTitle = 'This is a very long article title that should be displayed properly within the article card and should wrap correctly to maintain proper formatting and readability';
    
    const mockArticles = [
      { id: 'kb-001', title: longTitle }
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText(longTitle)).toBeInTheDocument();
  });

  it('should handle articles with special characters in titles', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'Article with special chars: @#$%^&*()' },
      { id: 'kb-002', title: 'Article with emojis: 🚀 💻 📱' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('Article with special chars: @#$%^&*()')).toBeInTheDocument();
    expect(screen.getByText('Article with emojis: 🚀 💻 📱')).toBeInTheDocument();
  });
});