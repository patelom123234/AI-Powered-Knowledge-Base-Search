import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';

import RelevantArticles from './RelevantArticles';

describe('RelevantArticles', () => {
    
  it('should render a list of article titles from props', () => {
    const mockArticles = [
      { id: 'kb-001', title: 'First Test Article' },
      { id: 'kb-002', title: 'Second Test Article' },
    ];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('First Test Article')).toBeInTheDocument();
    expect(screen.getByText('Second Test Article')).toBeInTheDocument();
  });

  it('should display a fallback message if no articles are provided', () => {
    const mockArticles = [];
    
    render(<RelevantArticles articles={mockArticles} />);

    expect(screen.getByText('No relevant articles found.')).toBeInTheDocument();
  });
});